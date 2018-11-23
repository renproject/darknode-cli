package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

// Providers have all the cloud service providers currently supported.
var Providers = []string{"aws", "do"}

// deployNode deploys node to the given cloud provider.
func deployNode(ctx *cli.Context) error {
	provider, err := provider(ctx)
	if err != nil {
		return err
	}

	switch provider {
	case "aws":
		return awsDeployment(ctx)
	case "do":
		return deployToDo(ctx)
	default:
		return ErrUnknownProvider
	}
}

// provider parses all provider flags and make sure only one provider is given.
func provider(ctx *cli.Context) (string, error) {
	var provider string

	counter := 0
	for i := range Providers {
		selected := ctx.Bool(Providers[i])
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

	fmt.Printf("%sDeploying dark nodes ... %s\n", GREEN, RESET)
	apply := fmt.Sprintf("cd %v && terraform apply -auto-approve", nodeDirectory)
	return run("bash", "-c", apply)
}

// outputUrl writes success message and the URL for registering the node
// to the terminal.
func outputUrl(nodeDir string) error {
	ip, err := getIp(nodeDir)
	if err != nil {
		return err
	}

	fmt.Printf("\n")
	fmt.Printf("%sCongratulations! Your Darknode is deployed.%s.\n", GREEN, RESET)
	fmt.Printf("%sJoin the network by registering your Darknode at%s\n", GREEN, RESET)
	fmt.Printf("%shttps://darknode.republicprotocol.com/status/%v%s\n", GREEN, ip, RESET)
	for i := 5; i >= 0; i-- {
		time.Sleep(time.Second)
		fmt.Printf("\r%sYou will be redirected to deregister your node in %v seconds%s", GREEN, i, RESET)
	}

	// Redirect the user to the registering URL.
	redirect, err := redirectCommand()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://darknode.republicprotocol.com/status/%v", ip)
	return run(redirect, url)
}
