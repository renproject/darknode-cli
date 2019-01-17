package main

import (
	"bufio"
	"context"
	"fmt"
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

	nodePath := nodePath(name)

	config, err := config.NewConfigFromJSONFile(path.Join(nodePath, "config.json"))
	if err != nil {
		return err
	}
	id := config.Address.ID()
	network := config.Ethereum.Network
	dnrAddress := common.HexToAddress(dnrAddress(network))
	ethereumNet := ethereumNetwork(network)

	// Query registry smart contract on Ethereum if the darknode is registered
	client, err := ethclient.Dial(fmt.Sprintf("https://%v.infura.io", ethereumNet))
	if err != nil {
		return err
	}
	registry, err := bindings.NewDarknodeRegistry(dnrAddress, client)
	if err != nil {
		return err
	}
	registered, err := registry.IsRegistered(&bind.CallOpts{}, common.BytesToAddress(id))
	if err != nil {
		return err
	}

	// Redirect the user to the de-registering URL if darknode is still registered.
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
		fmt.Printf("%sDeregister your darknode at %s.%s.\n", RED, url, RESET)

		for i := 5; i >= 0; i-- {
			time.Sleep(time.Second)
			fmt.Printf("\r%sYou will be redirected to deregister your node in %v seconds%s", RED, i, RESET)
		}
		fmt.Printf("%sPlease try again after you fully deregister your node%s\n", RED, RESET)

		redirect, err := redirectCommand()
		if err != nil {
			return err
		}

		return run(redirect, url)
	}

	// Check if the darknode is in pending deregistration state.
	pendingDeregistration, err := registry.IsPendingDeregistration(&bind.CallOpts{}, common.BytesToAddress(id))
	if err != nil {
		return err
	}
	if pendingDeregistration {
		fmt.Printf("%sYour node is pending for deregistration%s\n", RED, RESET)
		fmt.Printf("%sDarknode can only be destroyed when fully deregistred%s\n", RED, RESET)
		fmt.Printf("%sPlease wait for it to be fully deregistered and try again%s\n", RED, RESET)
		return nil
	}

	pendingRegistration, err := registry.IsPendingRegistration(&bind.CallOpts{}, common.BytesToAddress(id))
	if err != nil {
		return err
	}
	if pendingRegistration {
		fmt.Printf("%sYour node is pending for registration%s\n", RED, RESET)
		fmt.Printf("%sDarknode can only be destroyed when fully deregistred%s\n", RED, RESET)
		fmt.Printf("%sPlease deregister your node after the epoch shuffle and try again%s\n", RED, RESET)
		return nil
	}

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
	conn, err := contract.Connect(config.Ethereum)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(config.Keystore.EcdsaKey.PrivateKey)

	// Check if the darknode is refundable
	dnr, err := bindings.NewDarknodeRegistry(common.HexToAddress(conn.Config.DarknodeRegistryAddress), bind.ContractBackend(conn.Client))
	if err != nil {
		return err
	}
	refundable, err := dnr.IsRefundable(&bind.CallOpts{}, common.BytesToAddress(config.Address.ID()))
	if err != nil {
		return err
	}
	if !refundable {
		return fmt.Errorf("%sThe darknode is not refundable, please deregister your darknode first.%s\n", RED, RESET)
	}

	// Refund the bonds
	contractBinder, err := contract.NewBinder(auth, conn)
	if err != nil {
		return err
	}
	if err := contractBinder.Refund(config.Address.ID()); err != nil {
		if strings.Contains(err.Error(), "failed to estimate gas needed") {
			return ErrRejectedTx
		}
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
