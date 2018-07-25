package main

import (
"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

// deployNode deploys node to the given cloud provider.
func deployNode(ctx *cli.Context) error {
	provider, err := getProvider(ctx )
	if err != nil {
		return err
	}

	switch provider {
	case "aws":
		return deployToAWS(ctx)
	case "do":
		return deployToDigitalOcean(ctx)
	default:
		return ErrUnknownProvider
	}
}

// getProvider parses all provider flags and  make sure only
// one provider is present.
func getProvider(ctx *cli.Context) (string, error ) {
	var providers = []string{"aws", "do"}
	aws := ctx.Bool("aws")
	digitalOcean := ctx.Bool("do")

	// Make sure only one provider is provided
	counter, provider := 0, ""
	for i, j := range []bool{aws, digitalOcean} {
		if j {
			counter++
			provider = providers[i]
		}
	}

	if counter == 0 {
		return "", ErrNilProvider
	} else if counter > 1 {
		return "", ErrMultipleProviders
	}

	return provider, nil
}

// createNodeDirectory create the directory for the node.
func createNodeDirectory(ctx *cli.Context) (string , error) {
	name := ctx.String("name")
	tags := ctx.String("tags")
	nodeDir := nodeDirectory(name)

	// Make sure name is not nil
	if name == "" {
		return "", ErrEmptyNodeName
	}

	// Check if the directory exists or not.
	// todo : rm fail deployment if exists
	if _, err := os.Stat(nodeDir); !os.IsNotExist(err) {
		return "", ErrNodeExist
	}
	if err := os.Mkdir(nodeDir, 0777); err != nil {
		return "", err
	}

	// Store the tags
	if err := ioutil.WriteFile(nodeDir +"/tags.out", []byte(strings.TrimSpace(tags)), 0666); err != nil {
		return "", err
	}

	return name , nil
}


// runTerraform initializes and applies terraform
func runTerraform(nodeDirectory string) error {
	cmd := fmt.Sprintf("cd %v && terraform init", nodeDirectory)
	init := exec.Command("bash", "-c", cmd)
	pipeToStd(init)
	if err := init.Start(); err != nil {
		return err
	}
	if err := init.Wait(); err != nil {
		return err
	}

	fmt.Printf("%sDeploying dark nodes to AWS%s...\n", GREEN, RESET)

	cmd = fmt.Sprintf("cd %v && terraform apply -auto-approve", nodeDirectory)
	apply := exec.Command("bash", "-c", cmd)
	pipeToStd(apply)
	if err := apply.Start(); err != nil {
		return err
	}
	return apply.Wait()
}

// outputUrl writes success message and the URL for registering the node
// to the terminal.
func outputUrl (ctx *cli.Context, name, nodeDir string) error {
	network := ctx.String("network")
	ip, err := getIp(nodeDir)
	if err != nil {
		return err
	}

	// Update node to different branch according to the network.
	switch network {
	case "testnet":
	case "falcon":
		err = updateSingleNode(name, "develop", false)
	case "nightly":
		err = updateSingleNode(name, "nightly", false)
	}
	if err != nil {
		return err
	}

	fmt.Printf("\n")
	fmt.Printf("%sCongratulations! Your Darknode is deployed and running%s.\n", GREEN, RESET)
	fmt.Printf("%sJoin the network by registering your Darknode at%s\n", GREEN, RESET)
	fmt.Printf("%shttps://darknode.republicprotocol.com/status/%v%s\n", GREEN, ip, RESET)
	fmt.Printf("\n")
	return nil
}
