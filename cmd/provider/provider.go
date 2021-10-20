package provider

import (
	"encoding/json"
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
	DeployMultiple(ctx *cli.Context) error
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

// Provider returns the provider of a darknode instance.
func GetProvider(name string) (string, error) {
	if name == "" {
		return "", util.ErrEmptyName
	}

	cmd := fmt.Sprintf("cd %v && terraform output provider", util.NodePath(name))
	provider, err := util.CommandOutput(cmd)
	return strings.TrimSpace(provider), err
}

// initialise all files needed by deploying a new node
func initNode(name, tags string, network darknode.Network, configFile string) error {
	if err := initNodeDirectory(name, tags); err != nil {
		return err
	}
	if err := util.GenerateSshKeyAndWriteToDir(name); err != nil {
		return err
	}

	// Use given config for the new darknode
	var conf darknode.Config
	if configFile != "" {
		path, err := filepath.Abs(configFile)
		if err != nil {
			return errors.New("invalid config path")
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("cannot open config file, err = %v", err)
		}
		defer file.Close()

		if err := json.NewDecoder(file).Decode(&conf); err != nil {
			return err
		}
	} else {
		var err error
		conf, err = darknode.NewConfig(network)
		if err != nil {
			return err
		}
	}

	configData, err := json.MarshalIndent(conf, "", "	")
	if err != nil {
		return err
	}
	configPath := filepath.Join(util.NodePath(name), "config.json")
	return ioutil.WriteFile(configPath, configData, 0600)
}

func initNodeDirectory(name, tags string) error {
	if name == "" {
		return util.ErrEmptyName
	}
	path := util.NodePath(name)

	// Ask user to destroy the old node first if there's already a node with the name.
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("Node [%v] already exist. \nIf you want to do a fresh deployment, destroy the old one first.", name)
	}

	// Make a directory for the new node
	if err := os.Mkdir(path, 0700); err != nil {
		return err
	}

	// Create the `tags.out` file if not exist.
	tagsPath := filepath.Join(path, "tags.out")
	if _, err := os.Stat(tagsPath); err != nil {
		return ioutil.WriteFile(tagsPath, []byte(strings.TrimSpace(tags)), 0600)
	}

	return nil
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

func runTerraformSilent(name string) error {
	path := util.NodePath(name)
	init := fmt.Sprintf("cd %v && terraform init", path)
	if err := util.SilentRun("bash", "-c", init); err != nil {
		return err
	}

	apply := fmt.Sprintf("cd %v && terraform apply -auto-approve -no-color", path)
	return util.SilentRun("bash", "-c", apply)
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
