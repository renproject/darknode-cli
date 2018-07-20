package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

// destroyNode tears down the deployed darknode, but keep the config file.
func destroyNode(ctx *cli.Context) error {
	// FIXME : currently it only supports tear down AWS deployment.
	// Needs to figure out way which suits for all kinds of cloud service.
	force := ctx.Bool("force")
	name := ctx.String("name")

	if name == "" {
		cli.ShowCommandHelp(ctx, "down")
		return ErrEmptyNodeName
	}

	nodeDirectory := Directory + "/darknodes/" + name
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

	return destroyAwsNode(nodeDirectory)
}

// destroyAwsNode tears down the AWS instance.
func destroyAwsNode(nodeDirectory string) error {
	fmt.Printf("%sDestroying your darknode ...%s\n", GREEN, RESET)
	cmd := fmt.Sprintf("cd %v && terraform destroy --force && rm -rf %v", nodeDirectory, nodeDirectory)
	destroy := exec.Command("bash", "-c", cmd)
	pipeToStd(destroy)
	if err := destroy.Start(); err != nil {
		return err
	}

	return destroy.Wait()
}
