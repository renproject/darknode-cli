package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

// ErrNoDeploymentFound is returned when no node can be found for destroying
var ErrNoDeploymentFound = errors.New("cannot found any deployed node")

// ErrEmptyNodeName is returned when user doesn't provide the node name.
var ErrEmptyNodeName = errors.New("please provide the node name")

// destroyNode will tear down the deployed darknode, but keep the config file.
func destroyNode(ctx *cli.Context) error {
	// FIXME : currently it only supports tear down AWS deployment.
	// Needs to figure out way which suits for all kinds of cloud service.
	skip := ctx.Bool("skip")
	name := ctx.String("name")
	if name == "" {
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/" + name
	if !skip {
		ip, err := getIp(nodeDirectory)
		if err != nil {
			return ErrNoDeploymentFound
		}

		fmt.Printf("Make sure you have deregistered your node and withdrawn all fees.\n")
		fmt.Printf("You can do that by going to https://darknode.republicprotocol.com/ip4/%v\n", ip)
		fmt.Println("Have you deregistered your node and withdrawn all fees? (Yes/No)")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(text)) != "yes" {
			return nil
		}
	}

	return destroyAwsNode(nodeDirectory)
}

// destroyAwsNode tear down the AWS instance.
func destroyAwsNode(nodeDirectory string) error {
	log.Println("Destroying your darknode ...")
	cmd := fmt.Sprintf("cd %v && terraform destroy --force && rm -rf %v",nodeDirectory, nodeDirectory)
	destroy := exec.Command( "bash", "-c", cmd)
	pipeToStd(destroy)
	if err := destroy.Start(); err != nil {
		return err
	}
	return destroy.Wait()
}
