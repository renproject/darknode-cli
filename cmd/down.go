package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/contract/bindings"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/urfave/cli"
)

// destroyNode tears down the deployed darknode by its name.
func destroyNode(ctx *cli.Context) error {
	force := ctx.Bool("force")
	name := ctx.Args().First()
	if name == "" {
		cli.ShowCommandHelp(ctx, "down")
		return ErrEmptyNodeName
	}

	// Read the node config
	nodePath := nodePath(name)
	config, err := config.NewConfigFromJSONFile(path.Join(nodePath, "config.json"))
	if err != nil {
		return err
	}
	id := config.Address.ID()
	address := common.BytesToAddress(id)

	// Connect to Ethereum
	_, dnr, err := connect(config.Ethereum.Network)
	if err != nil {
		return err
	}
	// Check if the node is registered
	if err := checkRegistered(dnr, config.Ethereum.Network, id, address); err != nil {
		return err
	}
	// Check if the node is in pending registration/deregistration stage
	if err := checkPendingStage(dnr, address); err != nil {
		return err
	}
	// Check if the darknode has been refunded
	refunded, err := dnr.IsRefunded(&bind.CallOpts{}, address)
	if err != nil {
		return err
	}
	if !refunded {
		fmt.Printf("You haven't refund your darknode. Please refund with `darknode refund %v` command", name)
		return nil
	}

	// Check if user want to process without extra confirmation
	if !force {
		fmt.Println("Do you really want to destroy your darknode? (Yes/No)")

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		input := strings.ToLower(strings.TrimSpace(text))
		if input != "yes" && input != "y" {
			return nil
		}
	}
	fmt.Printf("%sDestroying your darknode ...%s\n", RESET, RESET)

	destroy := fmt.Sprintf("cd %v && terraform destroy --force && find . -type f -not -name 'config.json' -delete", nodePath)
	return run("bash", "-c", destroy)
}

// refund the REN bonds to the darknode operator.
func refund(ctx *cli.Context) error {
	name := ctx.Args().First()

	// Validate the name and check if the directory exists.
	nodePath, err := validateDarknodeName(name)
	if err != nil {
		return err
	}

	// Read the config and connect to Ethereum
	config, err := config.NewConfigFromJSONFile(nodePath + "/config.json")
	if err != nil {
		return err
	}
	client, dnr, err := connect(config.Ethereum.Network)
	if err != nil {
		return err
	}
	id := config.Address.ID()
	address := common.BytesToAddress(id)
	auth := bind.NewKeyedTransactor(config.Keystore.EcdsaKey.PrivateKey)

	refundable, err := dnr.IsRefundable(&bind.CallOpts{}, address)
	if err != nil {
		return err
	}
	if !refundable {
		// Check if the darknode has been refunded
		refunded, err := dnr.IsRefunded(&bind.CallOpts{}, address)
		if err != nil {
			return err
		}
		if refunded {
			return errors.New("you have already refunded the darknode")
		}

		// Check if the node is registered
		if err := checkRegistered(dnr, config.Ethereum.Network, id, address); err != nil {
			return err
		}
		// Check if the node is in pending registration/deregistration stage
		return checkPendingStage(dnr, address)
	}

	tx, err := dnr.Refund(auth, address)
	if err != nil {
		return err
	}
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	_, err = bind.WaitMined(timeoutCtx, client, tx)
	if err != nil {
		return err
	}
	fmt.Printf("%sYour REN bonds have been refunded to the operator account%s \n", GREEN, RESET)

	return nil
}

// Withdraw ETH and REN in the darknode address to the provided receiver address
func withdraw(ctx *cli.Context) error {
	name := ctx.Args().First()
	address := ctx.String("address")

	// Validate the name and received ethereum address
	nodePath, err := validateDarknodeName(name)
	if err != nil {
		return err
	}
	receiverAddr, err := stringToEthereumAddress(address)
	if err != nil {
		return err
	}

	// Read the darknode config
	config, err := config.NewConfigFromJSONFile(nodePath + "/config.json")
	if err != nil {
		return err
	}
	conn, err := contract.Connect(config.Ethereum)
	if err != nil {
		return err
	}
	darknodeEthAddress, err := republicAddressToEthAddress(config.Address.String())
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(config.Keystore.EcdsaKey.PrivateKey)
	auth.GasPrice = big.NewInt(5000000000) // Set GasPrise to 5 Gwei

	// Check REN balance first
	renAddress := renAddress(config.Ethereum.Network)
	if renAddress == "" {
		return ErrUnknownNetwork
	}
	tokenContract, err := bindings.NewERC20(common.HexToAddress(renAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		return err
	}
	renBalance, err := tokenContract.BalanceOf(&bind.CallOpts{}, darknodeEthAddress)
	if err != nil {
		return err
	}

	// Withdraw REN if the darknode has more than 1 REN.
	oneREN := big.NewInt(int64(math.Pow10(18)))
	if renBalance.Cmp(oneREN) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		tx, err := tokenContract.Transfer(auth, receiverAddr, renBalance)
		if err != nil {
			return err
		}
		receipt, err := conn.PatchedWaitMined(ctx, tx)
		if err != nil {
			return err
		}
		if receipt.Status == types.ReceiptStatusFailed {
			return ErrFailedTx
		}
		fmt.Printf("%sYour REN has been withdrawn from your darknode to [%v]%s \n", GREEN, receiverAddr.Hex(), RESET)
	}

	// Check ETH balance of the darknode
	balance, err := conn.Client.BalanceAt(context.Background(), darknodeEthAddress, nil)
	if err != nil {
		return err
	}
	transactionFee := big.NewInt(int64(5 * math.Pow10(9) * 21000)) //  5 Gwei Gas price
	// Transfer Eth back to the owner
	if balance.Cmp(transactionFee) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		value := new(big.Int)
		tx, err := conn.SendEth(ctx, auth, receiverAddr, value.Sub(balance, transactionFee))
		if err != nil {
			return err
		}
		receipt, err := conn.PatchedWaitMined(ctx, tx)
		if err != nil {
			return err
		}
		if receipt.Status == types.ReceiptStatusFailed {
			return ErrFailedTx
		}
		fmt.Printf("%sYour ETH has been withdrawn from your darknode to [%v]%s \n", GREEN, receiverAddr.Hex(), RESET)
	}

	return nil
}

// renAddress on different testnet
func renAddress(network contract.Network) string {
	switch network {
	case "mainnet":
		return "0x408e41876cCCDC0F92210600ef50372656052a38"
	case "testnet":
		return "0x6f429121a3bd3e6c1c17edbc676eec44cf117faf"
	default:
		return ""
	}
}

// darknode registry on different testnet
func dnrAddress(network contract.Network) string {
	switch network {
	case "mainnet":
		return "0x34bd421C7948Bc16f826Fd99f9B785929b121633"
	case "testnet":
		return "0x75Fa8349fc9C7C640A4e9F1A1496fBB95D2Dc3d5"
	default:
		return ""
	}
}

// ethereumNetwork returns the ethereum network name of the given republic
// protocol network
func ethereumNetwork(network contract.Network) string {
	switch network {
	case "mainnet":
		return "mainnet"
	case "testnet":
		return "kovan"
	default:
		return ""
	}
}

func connect(network contract.Network) (*ethclient.Client, *bindings.DarknodeRegistry, error) {
	dnrAddr := common.HexToAddress(dnrAddress(network))
	ethereumNet := ethereumNetwork(network)
	client, err := ethclient.Dial(fmt.Sprintf("https://%v.infura.io", ethereumNet))
	if err != nil {
		return nil, nil, err
	}
	dnr, err := bindings.NewDarknodeRegistry(dnrAddr, client)
	return client, dnr, err
}

func checkRegistered(dnr *bindings.DarknodeRegistry, network contract.Network, id identity.ID, address common.Address) error {
	registered, err := dnr.IsRegistered(&bind.CallOpts{}, address)
	if err != nil {
		return err
	}
	if registered {
		var url string
		switch network {
		case "testnet":
			url = fmt.Sprintf("https://dcc-testnet.republicprotocol.com/darknode/%v?action=deregister", id)
		case "mainnet":
			url = fmt.Sprintf("https://dcc.republicprotocol.com/darknode/%v?action=deregister", id)
		default:
			return ErrUnknownNetwork
		}

		fmt.Printf("%sYour node hasn't been deregistered%s\n", RED, RESET)
		fmt.Printf("%sDeregister your darknode at %s%s\n", RED, url, RESET)

		for i := 5; i >= 0; i-- {
			time.Sleep(time.Second)
			fmt.Printf("\r%sYou will be redirected to deregister your node in %v seconds%s", RED, i, RESET)
		}
		redirect, err := redirectCommand()
		if err != nil {
			return err
		}
		if err := run(redirect, url); err != nil {
			return err
		}
		return fmt.Errorf("%s\nPlease try again after you fully deregister your node%s\n", RED, RESET)
	}
	return nil
}

func checkPendingStage(dnr *bindings.DarknodeRegistry, address common.Address) error {
	pendingRegistration, err := dnr.IsPendingRegistration(&bind.CallOpts{}, address)
	if err != nil {
		return err
	}
	if pendingRegistration {
		return fmt.Errorf("%sYour node is currently in pending registration stage, please deregister your node after next epoch shuffle%s\n", RED, RESET)
	}
	pendingDeregistration, err := dnr.IsPendingDeregistration(&bind.CallOpts{}, address)
	if err != nil {
		return err
	}
	if pendingDeregistration {
		return fmt.Errorf("%sYour node is currently in pending deregistration stage, please wait for next epoch shuffle and try again%s\n", RED, RESET)
	}
	return nil
}
