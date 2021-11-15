package provider

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
)

var (
	// ErrUnknownProvider is returned when user tries to deploy a darknode with an unknown cloud provider.
	ErrUnknownProvider = errors.New("unknown cloud provider")

	// ErrUnsupportedInstanceType is returned when the selected instance type cannot be created to user account.
	ErrInstanceTypeNotAvailable = errors.New("selected instance type is not available")

	// ErrRegionNotAvailable is returned when the selected region is not available to user account.
	ErrRegionNotAvailable = errors.New("selected region is not available")
)

var (
	NameAws = "aws"
	NameDo  = "do"
	NameGcp = "gcp"
)

var darknodeService = `[Unit]
Description=RenVM Darknode Daemon
AssertPathExists=$HOME/.darknode

[Service]
WorkingDirectory=$HOME/.darknode
ExecStart=$HOME/.darknode/bin/darknode --config $HOME/.darknode/config.json
Restart=on-failure
PrivateTmp=true
NoNewPrivileges=true

# Specifies which signal to use when killing a service. Defaults to SIGTERM.
# SIGHUP gives parity time to exit cleanly before SIGKILL (default 90s)
KillSignal=SIGHUP

[Install]
WantedBy=default.target`

type Provider interface {
	Name() string
	Deploy(ctx *cli.Context) error
}

func ParseProvider(ctx *cli.Context) (Provider, error) {
	if ctx.Bool(NameAws) {
		return NewAws(ctx)
	}

	if ctx.Bool(NameDo) {
		return NewDo(ctx)
	}

	if ctx.Bool(NameGcp) {
		return NewGcp(ctx)
	}

	return nil, ErrUnknownProvider
}

// Validate the params which are general to all providers.
func validateCommonParams(ctx *cli.Context) error {
	// Check the name valida and not been used
	name := ctx.String("name")
	if err := util.ValidateName(name); err != nil {
		return err
	}
	if err := util.ValidateNodeExistence(name); err == nil {
		return fmt.Errorf("node [%v] already exist", name)
	}

	// Parse the network
	network := darknode.Network(ctx.String("network"))
	switch network {
	case darknode.Mainnet:
	case darknode.Testnet:
	case darknode.Devnet:
	default:
		return errors.New("unknown RenVM network")
	}

	// Verify the config file if user wants to use their own config
	configFile := ctx.String("config")
	if configFile != "" {
		// verify the config exist and of the right format
		path, err := filepath.Abs(configFile)
		if err != nil {
			return err
		}
		if _, err := os.Stat(path); err != nil {
			return errors.New("config file doesn't exist")
		}
		_, err = darknode.NewConfigFromJSONFile(path)
		if err != nil {
			return fmt.Errorf("incompatible config, err = %w", err)
		}
	}
	return nil
}

// Provider returns the provider of a darknode instance.
func GetProvider(name string) (string, error) {
	if name == "" {
		return "", util.ErrEmptyName
	}

	cmd := fmt.Sprintf("cd %v && terraform output provider", util.NodePath(name))
	provider, err := util.CommandOutput(cmd)
	if strings.HasPrefix(provider, "\"") {
		provider = strings.Trim(provider, "\"")
	}
	return strings.TrimSpace(provider), err
}

// initialise all files needed by deploying a new node
func initNode(ctx *cli.Context) error {
	name := ctx.String("name")
	path := util.NodePath(name)
	configFile := ctx.String("config")
	network := darknode.Network(ctx.String("network"))

	// Create directory for the Darknode
	if err := os.MkdirAll(path, 0700); err != nil {
		return err
	}

	// Create `tags.out` file
	tags := []byte(strings.TrimSpace(ctx.String("tags")))
	tagsPath := filepath.Join(path, "tags.out")
	if err := ioutil.WriteFile(tagsPath, tags, 0600); err != nil {
		return err
	}

	// Create `ssh_keypair` and `ssh_keypair.pub` files for the remote instance
	if err := util.GenerateSshKeyAndWriteToDir(name); err != nil {
		return err
	}

	// Use given config for the new darknode
	var conf darknode.Config
	if configFile != "" {
		var err error
		conf, err = darknode.NewConfigFromJSONFile(path)
		if err != nil {
			return errors.New("invalid config file")
		}
	} else {
		var err error
		conf, err = darknode.NewConfig(network)
		if err != nil {
			return err
		}
	}

	return darknode.ConfigToFile(conf, filepath.Join(util.NodePath(name), "config.json"))
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
func outputURL(name string) error {
	url, err := util.RegisterUrl(name)
	if err != nil {
		return err
	}
	color.Green("")
	color.Green("Congratulations! Your Darknode is deployed.")
	color.Green("Join the network by registering your Darknode at %s", url)
	return util.OpenInBrowser(url)
}
