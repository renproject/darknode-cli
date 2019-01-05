package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

// Provider represents a Terraform provider
type Provider string

const (
	AWS           Provider = "aws"
	DIGITAL_OCEAN Provider = "do"
)

// Providers have all the cloud service providers currently supported.
var Providers = []Provider{AWS, DIGITAL_OCEAN}

// deployNode deploys node to the given cloud provider.
func deployNode(ctx *cli.Context) error {
	provider, err := provider(ctx)
	if err != nil {
		return err
	}

	switch provider {
	case AWS:
		return awsDeployment(ctx)
	case DIGITAL_OCEAN:
		return deployToDo(ctx)
	default:
		return ErrUnknownProvider
	}
}

// provider parses all provider flags and make sure only one provider is given.
func provider(ctx *cli.Context) (Provider, error) {
	var provider Provider

	counter := 0
	for i := range Providers {
		selected := ctx.Bool(string(Providers[i]))
		if selected {
			counter++
			provider = Providers[i]
		}
	}

	switch counter {
	case 0:
		return "", ErrNilProvider
	case 1:
		return provider, nil
	default:
		return "", ErrMultipleProviders
	}
}

// mkdir creates the directory for the darknode.
func mkdir(name, tags string) error {
	if name == "" {
		return ErrEmptyNodeName
	}
	nodePath := nodePath(name)

	// Check if the directory exists or not.
	if _, err := os.Stat(nodePath); err == nil {
		if _, err := os.Stat(nodePath + "/multiAddress.out"); os.IsNotExist(err) {
			// todo : need to ask user whether they want to use the old config.
			if err := run("rm", "-rf", nodePath); err != nil {
				return err
			}
		} else {
			return ErrNodeExist
		}
	}
	if err := os.Mkdir(nodePath, 0777); err != nil {
		return err
	}

	// Store the tags
	return ioutil.WriteFile(nodePath+"/tags.out", []byte(strings.TrimSpace(tags)), 0666)
}

// runTerraform initializes and applies terraform
func runTerraform(nodeDirectory string) error {
	init := fmt.Sprintf("cd %v && terraform init", nodeDirectory)
	if err := run("bash", "-c", init); err != nil {
		return err
	}

	fmt.Printf("\n%sDeploying dark nodes ... %s\n", RESET, RESET)
	apply := fmt.Sprintf("cd %v && terraform apply -auto-approve -no-color", nodeDirectory)
	return run("bash", "-c", apply)
}

// outputURL writes success message and the URL for registering the node
// to the terminal.
func outputURL(nodeDir, name, network string, publicKey []byte) error {
	id, err := getID(nodeDir)
	if err != nil {
		return err
	}

	publicKeyHex := hex.EncodeToString(publicKey)
	var url string
	switch network {
	case "mainnet":
		url = fmt.Sprintf("https://darknode-center-mainnet.herokuapp.com/darknode/%s?action=register&public_key=0x%s&name=%s", id, publicKeyHex, name)
	case "testnet":
		url = fmt.Sprintf("https://darknode-center-testnet.herokuapp.com/darknode/%s?action=register&public_key=0x%s&name=%s", id, publicKeyHex, name)
	}

	fmt.Printf("\n")
	fmt.Printf("%sCongratulations! Your Darknode is deployed.%s\n\n", GREEN, RESET)
	fmt.Printf("%sJoin the network by registering your Darknode at %s%s\n\n", GREEN, url, RESET)
	for i := 5; i > 0; i-- {
		fmt.Printf("\rYou will be redirected to register your node in %v seconds", i)
		time.Sleep(time.Second)
	}
	fmt.Printf("\r")

	// Redirect the user to the registering URL.
	redirect, err := redirectCommand()
	if err != nil {
		return err
	}
	return run(redirect, url)
}
