package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/darknode/bindings"
	"github.com/renproject/darknode-cli/util"
	"github.com/renproject/mercury/sdk/client/ethclient"
	"github.com/renproject/mercury/types/ethtypes"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// status represents the registration status of a Darknode.
type status int

const (
	nilStatus status = iota // Either not registered or fully deregistered.
	pendingRegistration
	registered
	pendingDeregistration
	notRefunded
)

// err returns the error message of invalid status.
func (s status) err() string {
	switch s {
	case pendingRegistration:
		return "Darknode is currently pending registration."
	case registered:
		return "Darknode is still registered."
	case pendingDeregistration:
		return "Darknode is currently pending deregistration."
	case notRefunded:
		return "Darknode bond has not been withdrawn."
	default:
		return ""
	}
}

// destroyNode tears down the deployed darknode by its name.
func destroyNode(ctx *cli.Context) error {
	force := ctx.Bool("force")
	name := ctx.Args().First()
	path := util.NodePath(name)

	// Check node current registration status.
	if !force {
		st, err := nodeStatus(name)
		if err != nil {
			color.Red("Failed to get Darknode registration status: %v", err)
		}
		switch st {
		case pendingRegistration, pendingDeregistration, registered, notRefunded:
			color.Red(st.err())
			color.Red("Please try again once your Darknode has been fully deregistered and refunded.")
			return nil
		default:
		}

		// Last time confirm with user.
		fmt.Println("Do you really want to destroy your darknode? (Yes/No)")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		input := strings.ToLower(strings.TrimSpace(text))
		if input != "yes" && input != "y" {
			return nil
		}
	}

	color.Green("Backing up config...")
	if err := util.BackUpConfig(name); err != nil {
		return err
	}

	color.Green("Destroying your Darknode...")
	destroy := fmt.Sprintf("cd %v && terraform destroy --force && cd .. && rm -rf %v", path, name)
	return util.Run("bash", "-c", destroy)
}

// Withdraw ETH and REN in the darknode address to the provided receiver address
func withdraw(ctx *cli.Context) error {
	name := ctx.Args().First()
	withdrawAddress := ctx.String("address")

	// Validate the name and received ethereum address
	if !common.IsHexAddress(withdrawAddress) {
		return errors.New("invalid receiver address")
	}
	receiverAddr := common.HexToAddress(withdrawAddress)

	// Parse the node config
	config, err := util.Config(name)
	if err != nil {
		return err
	}

	// Connect to Ethereum
	client, err := connect(config.Network)
	if err != nil {
		return err
	}

	// Create a transactor for ethereum tx
	ethAddr := crypto.PubkeyToAddress(config.Keystore.Ecdsa.PublicKey)
	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	auth := bind.NewKeyedTransactor(config.Keystore.Ecdsa.PrivateKey)
	auth.GasPrice = big.NewInt(5000000000) // Set GasPrise to 5 Gwei
	auth.Context = c

	// Check REN balance first
	renAddress := renAddress(config.Network)
	tokenContract, err := bindings.NewERC20(common.HexToAddress(renAddress), client.EthClient())
	if err != nil {
		return err
	}
	renBalance, err := tokenContract.BalanceOf(&bind.CallOpts{}, ethAddr)
	if err != nil {
		return err
	}

	// Withdraw REN if the darknode has more than 1 REN.
	fmt.Println("Checking REN balance...")
	oneREN := big.NewInt(1e18)
	if renBalance.Cmp(oneREN) > 0 {
		tx, err := tokenContract.Transfer(auth, receiverAddr, renBalance)
		if err != nil {
			return err
		}
		receipt, err := bind.WaitMined(c, client.EthClient(), tx)
		if err != nil {
			return err
		}
		renBalanceNoDecimals := big.NewInt(0).Div(renBalance, oneREN)
		color.Green("%v REN has been withdrawn from your darknode to [%v]. TxHash: %v.", renBalanceNoDecimals.Int64(), receiverAddr.Hex(), receipt.TxHash.Hex())
	} else {
		color.Green("Your account doesn't have REN token.")
	}

	// Check the ETH balance
	fmt.Println("Checking ETH balance...")
	balance, err := client.Balance(c, ethtypes.Address(ethAddr))
	if err != nil {
		return err
	}
	gas := ethtypes.Gwei(5 * 21000)
	zero := ethtypes.Wei(0)
	if balance.Gt(zero) {
		if balance.Gt(gas) {
			tx, err := transfer(auth, receiverAddr, balance.Sub(gas), client)
			if err != nil {
				return err
			}
			color.Green("Your ETH has been withdrawn from your darknode to [%v]. TxHash: %v.", receiverAddr.Hex(), tx.Hash().Hex())
		} else {
			return fmt.Errorf("your account has %v wei which is not enough to cover the transaction fee %v on ethereum", balance, gas)
		}
	} else {
		color.Green("Your don't have any ETH left in your account.")
	}
	return nil
}

// transfer ETH to the provided address.
func transfer(transactor *bind.TransactOpts, receiver common.Address, amount ethtypes.Amount, client ethclient.Client) (*types.Transaction, error) {
	bound := bind.NewBoundContract(receiver, abi.ABI{}, nil, client.EthClient(), nil)
	transactor.Value = amount.ToBig()
	transactor.GasLimit = 21000
	return bound.Transfer(transactor)
}

// renAddress on different network
func renAddress(network darknode.Network) string {
	switch network {
	case darknode.Mainnet, darknode.Chaosnet:
		return "0x408e41876cCCDC0F92210600ef50372656052a38"
	case darknode.Testnet, darknode.Devnet:
		return "0x2CD647668494c1B15743AB283A0f980d90a87394"
	default:
		panic("unknown network")
	}
}

// connect to Ethereum.
func connect(network darknode.Network) (ethclient.Client, error) {
	logger := logrus.New()
	switch network {
	case darknode.Mainnet, darknode.Chaosnet:
		return ethclient.New(logger, ethtypes.Mainnet)
	case darknode.Testnet, darknode.Devnet:
		return ethclient.New(logger, ethtypes.Kovan)
	default:
		return nil, errors.New("unknown network")
	}
}

// nodeStatus returns the registration status of the darknode with given name.
func nodeStatus(name string) (status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	config, err := util.Config(name)
	if err != nil {
		return 0, err
	}
	address := crypto.PubkeyToAddress(config.Keystore.Ecdsa.PublicKey)

	// Connect to Ethereum
	client, err := connect(config.Network)
	if err != nil {
		return 0, err
	}
	dnrAddr, err := config.DnrAddr(client.EthClient())
	if err != nil {
		return 0, err
	}
	dnr, err := bindings.NewDarknodeRegistry(dnrAddr, client.EthClient())
	if err != nil {
		return 0, err
	}

	// Check if node is in pending registration status
	pr, err := dnr.IsPendingRegistration(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return 0, err
	}
	if pr {
		return pendingRegistration, nil
	}

	// Check if node is registered
	r, err := dnr.IsRegistered(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return 0, err
	}
	if r {
		return registered, nil
	}

	// Check if node in pending deregistration status
	pd, err := dnr.IsPendingDeregistration(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return 0, err
	}
	if pd {
		return pendingDeregistration, nil
	}

	// Check if node has been refunded
	refunded, err := dnr.IsRefunded(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return 0, err
	}
	if !refunded {
		return notRefunded, nil
	}
	return nilStatus, nil
}
