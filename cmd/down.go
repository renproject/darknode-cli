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

func refund(ctx *cli.Context) error {
	name := ctx.Args().First()
	refundAll := ctx.Bool("all")

	// Validate the name and check if the directory exists.
	if name == "" {
		return ErrEmptyNodeName
	}
	nodeDir := nodeDirectory(name)
	if _, err := os.Stat(nodeDir); err != nil {
		return ErrNodeNotExist
	}
	if _, err := os.Stat(nodeDir + "/config.json"); os.IsNotExist(err) {
		return ErrNodeNotExist
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
	auth.GasPrice = big.NewInt(10000000000) // Set GasPrise to 10 GWEI
	contractBinder, err := contract.NewBinder(auth, conn)
	if err != nil {
		return err
	}
	// TODO : only the darknode owner can call the refund function, disable for now
	// TODO : and wait for the result of contract auditing.
	// if err := contractBinder.Refund(config.Address.ID()); err != nil {
	// 	return err
	// }
	// fmt.Printf("%sYour REN bonds have been refunded to your nominated address%s \n", GREEN, RESET)

	// Refund ETH and REN in the darknode address if there are
	if refundAll {
		ownerAddress, err := contractBinder.GetOwner(config.Address.ID())
		if err != nil {
			return err
		}
		darknodeEthAddress, err := republicAddressToEthAddress(config.Address.String())
		if err != nil {
			return err
		}
		balance, err := conn.Client.BalanceAt(context.Background(), darknodeEthAddress, nil)
		if err != nil {
			return err
		}
		transactionFee := big.NewInt(int64(math.Pow10(15))) // 0.001 ETH as tx fee

		// Transfer Eth back to the owner
		if balance.Cmp(transactionFee) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			value := new(big.Int)
			tx, err := conn.SendEth(ctx, auth, ownerAddress, value.Sub(balance, transactionFee))
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

		// Transfer REN back to the owners if they accidentally send REN to the darknodes.
		renAddress := renAddress(config.Ethereum.Network)
		if renAddress == "" {
			return ErrUnknownNetwork
		}
		tokenContract, err := bindings.NewERC20(common.HexToAddress(renAddress), bind.ContractBackend(conn.Client))
		balance, err = tokenContract.BalanceOf(&bind.CallOpts{}, darknodeEthAddress)
		if err != nil {
			return err
		}
		oneREN := big.NewInt(int64(math.Pow10(18))) // 0.001 ETH as tx fee
		// Transfer REN back to the owner
		if balance.Cmp(oneREN) > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			tx, err := tokenContract.Transfer(auth, ownerAddress, balance)
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
			fmt.Printf("%sWe have refund the REN in your darknode address%s \n", GREEN, RESET)
		}
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
