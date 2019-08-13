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
	"github.com/renproject/libeth-go"
	"github.com/republicprotocol/darknode-cli/bindings"
	"github.com/republicprotocol/darknode-cli/config"
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
	network := config.Ethereum.Network
	ethAddr, err := config.EthereumAdress()
	if err != nil {
		return err
	}

	// Connect to Ethereum
	_, dnr, err := connect(network)
	if err != nil {
		return err
	}

	// Check if the node is registered
	if err := checkRegistered(dnr, network, config); err != nil {
		return err
	}
	// Check if the node is in pending registration/deregistration stage
	if err := checkPendingStage(dnr, config); err != nil {
		return err
	}

	// Check if the darknode has been refunded
	context, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	refunded, err := dnr.IsRefunded(&bind.CallOpts{Context: context}, ethAddr)
	if err != nil {
		return err
	}
	if !refunded {
		fmt.Println("You haven't refund your darknode. Please refund your darknode from the command center")
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

// Withdraw ETH and REN in the darknode address to the provided receiver address
func withdraw(ctx *cli.Context) error {
	name := ctx.Args().First()
	withdrawAddress := ctx.String("address")

	// Validate the name and received ethereum address
	nodePath, err := validateDarknodeName(name)
	if err != nil {
		return err
	}
	receiverAddr, err := stringToEthereumAddress(withdrawAddress)
	if err != nil {
		return err
	}

	// Read config of the specified darknode
	config, err := config.NewConfigFromJSONFile(path.Join(nodePath, "config.json"))
	if err != nil {
		return err
	}
	network := config.Ethereum.Network

	// Connect to Ethereum
	client, _, err := connect(network)
	if err != nil {
		return err
	}
	account, err := libeth.NewAccount(client, config.Keystore.Ecdsa.PrivateKey)
	if err != nil {
		return err
	}

	darknodeEthAddress, err := republicAddressToEthAddress(config.Address)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(config.Keystore.Ecdsa.PrivateKey)
	auth.GasPrice = big.NewInt(5000000000) // Set GasPrise to 5 Gwei

	// Check REN balance first
	renAddress := renAddress(network)
	if renAddress == "" {
		return ErrUnknownNetwork
	}
	tokenContract, err := bindings.NewERC20(common.HexToAddress(renAddress), client.EthClient())
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
		tx, err := tokenContract.Transfer(auth, receiverAddr, renBalance)
		if err != nil {
			return err
		}
		minedCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_, err = bind.WaitMined(minedCtx, client.EthClient(), tx)
		if err != nil {
			return err
		}
		fmt.Printf("%sYour REN has been withdrawn from your darknode to [%v]. TxHash: %v.%s\n", GREEN, receiverAddr.Hex(), tx.Hash().Hex(), RESET)
	}

	ethCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Check the ETH balance
	balance, err := account.BalanceAt(ethCtx, nil)
	if err != nil {
		return err
	}
	if balance.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	tx, err := account.Transfer(ethCtx, receiverAddr, nil, libeth.Fast, 0, true)
	if err != nil {
		return err
	}
	fmt.Printf("%sYour ETH has been withdrawn from your darknode to [%v]. TxHash: %v.%s\n", GREEN, receiverAddr.Hex(), tx.Hash().Hex(), RESET)
	return nil
}

// renAddress on different network
func renAddress(network string) string {
	switch network {
	case "mainnet":
		return "0x408e41876cCCDC0F92210600ef50372656052a38"
	case "kovan", "testnet":
		return "0x2CD647668494c1B15743AB283A0f980d90a87394"
	default:
		return ""
	}
}

// renAddress on different network
func dnrAddress(network string) common.Address {
	switch network {
	case "mainnet":
		return common.HexToAddress("0x34bd421C7948Bc16f826Fd99f9B785929b121633")
	case "kovan", "testnet":
		return common.HexToAddress("0x75Fa8349fc9C7C640A4e9F1A1496fBB95D2Dc3d5")
	default:
		panic(fmt.Sprintf("unknown network %v", network))
	}
}

func connect(network string) (libeth.Client, *bindings.DarknodeRegistry, error) {
	client, err := libeth.NewMercuryClient(network, "dcc")
	if err != nil {
		return libeth.Client{}, nil, err
	}
	contract, err := bindings.NewDarknodeRegistry(dnrAddress(network), client.EthClient())
	if err != nil {
		return libeth.Client{}, nil, err
	}

	return client, contract, nil
}

func checkRegistered(dnr *bindings.DarknodeRegistry, network string, config config.Config) error {
	ethAddr, err := config.EthereumAdress()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	registered, err := dnr.IsRegistered(&bind.CallOpts{Context: ctx}, ethAddr)
	if err != nil {
		return err
	}
	if registered {
		var url string
		switch network {
		case "testnet", "kovan":
			url = fmt.Sprintf("https://dcc-testnet.republicprotocol.com/darknode/%v?action=deregister", config.Address)
		case "mainnet":
			url = fmt.Sprintf("https://dcc.republicprotocol.com/darknode/%v?action=deregister", config.Address)
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

func checkPendingStage(dnr *bindings.DarknodeRegistry, config config.Config) error {
	ethAddr, err := config.EthereumAdress()
	if err != nil {
		return err
	}

	reCtx, reCancel := context.WithTimeout(context.Background(), time.Minute)
	defer reCancel()
	pendingRegistration, err := dnr.IsPendingRegistration(&bind.CallOpts{Context: reCtx}, ethAddr)
	if err != nil {
		return err
	}
	if pendingRegistration {
		return fmt.Errorf("%sYour node is currently in pending registration stage, please deregister your node after next epoch shuffle%s\n", RED, RESET)
	}

	deCtx, deCancel := context.WithTimeout(context.Background(), time.Minute)
	defer deCancel()
	pendingDeregistration, err := dnr.IsPendingDeregistration(&bind.CallOpts{Context: deCtx}, ethAddr)
	if err != nil {
		return err
	}
	if pendingDeregistration {
		return fmt.Errorf("%sYour node is currently in pending deregistration stage, please wait for next epoch shuffle and try again%s\n", RED, RESET)
	}

	return nil
}
