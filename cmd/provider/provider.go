package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/republicprotocol/darknode-cli/darknode"
	"github.com/republicprotocol/darknode-cli/util"
	"github.com/urfave/cli"
)

// ErrEmptyNodeName is returned when user doesn't provide the node name.
var ErrEmptyNodeName = errors.New("node name cannot be empty")

// ErrUnknownProvider is returned when user tries to deploy a darknode with an unknown cloud provider.
var ErrUnknownProvider = errors.New("unknown cloud provider")

// ErrUnsupportedInstanceType is returned when the selected instance type cannot be created to user account.
var ErrInstanceTypeNotAvailable = errors.New("selected instance type is not available")

// ErrRegionNotAvailable is returned when the selected region is not available to user account.
var ErrRegionNotAvailable = errors.New("selected region is not available")

var (
	NameAws = "aws"
	NameDo  = "do"
	NameGcp = "gcp"
)

type Provider interface {
	Name() string
	Deploy(ctx *cli.Context) error
}

func ParseProvider(ctx *cli.Context) (Provider, error) {
	if ctx.Bool(NameAws) {
		return NewAws(ctx)
	}

	// if ctx.Bool(NameDo) {
	// 	return
	// }
	//
	// if ctx.Bool(NameGcp) {
	// 	return
	// }

	return nil, ErrUnknownProvider
}

// Provider returns the provider of a darknode instance.
func GetProvider(name string) (string, error) {
	// Get main.tf file
	filePath := filepath.Join(util.NodePath(name), "main.tf")
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Check if it's aws or digital ocean
	if strings.Contains(string(data), `provider "aws"`) {
		return NameAws, nil
	} else if strings.Contains(string(data), `provider "digitalocean"`) {
		return NameDo, nil
	} else if strings.Contains(string(data), `provider "gcp"`) {
		return NameGcp, nil
	} else {
		return "", errors.New("unknown provider")
	}
}

// initialise all files needed by deploying a new node
func initNode(name, tags string, network darknode.Network) error {
	if err := util.InitNodeDirectory(name, tags); err != nil {
		return err
	}
	if err := util.NewKey(name); err != nil {
		return err
	}

	// Generate a new config and write to a file.
	config, err := darknode.NewConfig(network)
	if err != nil {
		return err
	}
	configData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	configPath := filepath.Join(util.NodePath(name), "config.json")
	return ioutil.WriteFile(configPath, configData, 0600)
}

func runTerraform(name string) error {
	path := util.NodePath(name)
	init := fmt.Sprintf("cd %v && terraform init", path)
	if err := util.Run("bash", "-c", init); err != nil {
		return err
	}

	fmt.Println("Deploying darknode ... ")
	apply := fmt.Sprintf("cd %v && terraform apply -auto-approve -no-color", path)
	return util.Run("bash", "-c", apply)
}

// outputURL writes success message and the URL for registering the node to the terminal.
func outputURL(name string, network darknode.Network) error {
	// id, err := util.ID(name)
	// if err != nil {
	// 	return err
	// }

	// FIXME: what kind of encoding we want for the public key
	// publicKeyHex := hex.EncodeToString(publicKey)

	// TODO : Print the DCC url for user to register their darknode
	var url string
	switch network {
	case darknode.Mainnet:
	case darknode.Chaosnet:
	case darknode.Testnet:
	case darknode.Devnet:
	default:
		panic("unknown network")
	}
	url = "https://www.renproject.io"

	fmt.Printf("\n")
	fmt.Printf("%sCongratulations! Your Darknode is deployed.%s\n\n", util.GREEN, util.RESET)
	fmt.Printf("%sJoin the network by registering your Darknode at %s%s\n\n", util.GREEN, url, util.RESET)
	return nil
}
