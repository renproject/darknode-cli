package main

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/contract/bindings"
	"github.com/urfave/cli"
)

// destroyNode tears down the deployed darknode by its name.
func destroyNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	force := ctx.Bool("force")

	if name == "" {
		cli.ShowCommandHelp(ctx, "down")
		return ErrEmptyNodeName
	}

	nodeDirectory := nodeDirectory(name)
	if !force {
		ip, err := getIp(nodeDirectory)
		if err != nil {
			return ErrNoDeploymentFound
		}

		for {
			fmt.Printf("You need to %sderegister your Darknode%s and %swithdraw all fees%s at\n", RED, RESET, RED, RESET)
			fmt.Printf("https://darknode.republicprotocol.com/status/%v\n", ip)
			fmt.Println("Have you deregistered your Darknode and withdrawn all fees? (Yes/No)")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			input := strings.ToLower(strings.TrimSpace(text))
			if input == "yes" || input == "y" {
				break
			}
			if input == "no" || input == "n" {
				return nil
			}
		}
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
		return err
	}
	fmt.Printf("%sYour REN bonds have been refunded to your nominated address%s \n", GREEN, RESET)

	return nil
}

// Withdraw ETH and REN in the darknode address to the given address
func withdraw(ctx *cli.Context) error {
	name := ctx.Args().First()
	address := ctx.String("address")

	// Validate the name and Check if the node exists
	nodeDir, err := validateDarknodeName(name)
	if err != nil {
		return err
	}

	// Validate the receiver ethereum address
	if address == "" {
		return ErrEmptyAddress
	}
	if !common.IsHexAddress(address) {
		return ErrInvalidEthereumAddress
	}
	receiverAddr := common.HexToAddress(address)

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
	auth.GasPrice = big.NewInt(5000000000) // Set GasPrise to 5 GWEI

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
		fmt.Printf("%sWe have withdraw all the REN in your darknode address%s \n", GREEN, RESET)
	}

	// Check darknode balance
	balance, err := conn.Client.BalanceAt(context.Background(), darknodeEthAddress, nil)
	if err != nil {
		return err
	}
	transactionFee := big.NewInt(int64(5 * math.Pow10(9) * 21001)) //  5 Gwei Gas price

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
		fmt.Printf("%sWe have refund the ETH in your darknode address%s \n", GREEN, RESET)
	}

	return nil
}

// renAddress on different testnet
func renAddress(network contract.Network) string {
	switch network {
	case "testnet":
		return "0x6f429121a3bd3e6c1c17edbc676eec44cf117faf"
	case "falcon":
		return "0x87e83f957a2f3a2e5fe16d5c6b22e38fd28bdc06"
	case "nightly":
		return "0x15f692d6b9ba8cec643c7d16909e8acdec431bf6"
	default:
		return ""
	}
}
