package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	dnr "github.com/republicprotocol/darknode-cli/bindings"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/contract/bindings"
	"github.com/urfave/cli"
)

// destroyNode tears down the deployed darknode by its name.
func destroyNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	if name == "" {
		cli.ShowCommandHelp(ctx, "down")
		return ErrEmptyNodeName
	}

	nodeDirectory := nodeDirectory(name)
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return ErrNoDeploymentFound
	}

	config, err := config.NewConfigFromJSONFile(path.Join(nodeDirectory, "config.json"))
	if err != nil {
		return err
	}
	id := config.Address.ID()
	network := config.Ethereum.Network
	dnrAddress := common.HexToAddress(dnrAddress(network))
	testnet := ethereumTestnet(network)

	// Query registry smart contract on Ethereum if the darknode is registered
	client, err := ethclient.Dial(fmt.Sprintf("https://%v.infura.io", testnet))
	if err != nil {
		return err
	}
	registry, err := dnr.NewBindings(dnrAddress, client)
	if err != nil {
		return err
	}
	registered, err := registry.IsRegistered(&bind.CallOpts{}, common.BytesToAddress(id))
	if err != nil {
		return err
	}

	// Redirect the user to the de-registering URL if darknode is still registered.
	if registered {
		fmt.Printf("%sYour node hasn't been deregistered%s\n", RED, RESET)
		fmt.Printf("%sYou will be redirected to deregister your node%s\n", RED, RESET)
		time.Sleep(3 * time.Second)

		var redirect *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			redirect = exec.Command("open", fmt.Sprintf("https://darknode.republicprotocol.com/status/%v", ip))
		case "linux":
			redirect = exec.Command("xdg-open", fmt.Sprintf("https://darknode.republicprotocol.com/status/%v", ip))
		default:
			return errors.New("unsupported operating system")
		}
		pipeToStd(redirect)
		if err := redirect.Start(); err != nil {
			return err
		}
		return redirect.Wait()
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

	fmt.Printf("%sDestroying your darknode ...%s\n", GREEN, RESET)
	cmd := fmt.Sprintf("cd %v && terraform destroy --force && find . -type f -not -name 'config.json' -delete", nodeDirectory)
	destroy := exec.Command("bash", "-c", cmd)
	pipeToStd(destroy)
	if err := destroy.Start(); err != nil {
		return err
	}

	return destroy.Wait()
}

// refund the REN bonds to the darknode operator.
func refund(ctx *cli.Context) error {
	name := ctx.Args().First()

	// Validate the name and check if the directory exists.
	nodeDir, err := validateDarknodeName(name)
	if err != nil {
		return err
	}

	// Read the config and refund the REN bonds
	config, err := config.NewConfigFromJSONFile(nodeDir + "/config.json")
	if err != nil {
		return err
	}
	conn, err := contract.Connect(config.Ethereum)
	if err != nil {
		return err
	}
	auth := bind.NewKeyedTransactor(config.Keystore.EcdsaKey.PrivateKey)
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
	nodeDir, err := validateDarknodeName(name)
	if err != nil {
		return err
	}
	receiverAddr, err := stringToEthereumAddress(address)
	if err != nil {
		return err
	}

	// Read the darknode config
	config, err := config.NewConfigFromJSONFile(nodeDir + "/config.json")
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
		fmt.Printf("%sAll the REN in your darknode address have been withdrawed to [%v]%s \n", GREEN, receiverAddr.Hex(), RESET)
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
		fmt.Printf("%sAll the ETH in your darknode address have been withdrawed to [%v]%s \n", GREEN, receiverAddr.Hex(), RESET)
	}

	return nil
}

// renAddress on different testnet
func renAddress(network contract.Network) string {
	switch network {
	case "mainnet":
		return "0x81793734c6Cf6961B5D0D2d8a30dD7DF1E1803f1"
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
		return "0x3799006a87FDE3CCFC7666B3E6553B03ED341c2F"
	case "testnet":
		return "0x75Fa8349fc9C7C640A4e9F1A1496fBB95D2Dc3d5"
	default:
		return ""
	}
}

// ethereumTestnet returns the testnet name of different network
func ethereumTestnet(network contract.Network) string {
	switch network {
	case "mainnet":
		return "mainnet"
	case "testnet":
		return "kovan"
	default:
		return ""
	}
}
