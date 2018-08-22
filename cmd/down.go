package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

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
	// todo : keep the keystore and config file
	cmd := fmt.Sprintf("cd %v && terraform destroy --force && rm -rf %v", nodeDirectory, nodeDirectory)
	destroy := exec.Command("bash", "-c", cmd)
	pipeToStd(destroy)
	if err := destroy.Start(); err != nil {
		return err
	}

	return destroy.Wait()
}

func refund(ctx *cli.Context) error {
	name := ctx.Args().First()
	operator :=  ctx.Args().Get(1)
	all := ctx.Bool("all")

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


	// Validate the Ethereum address
	ethAddressRegex := "(0x)?[a-fA-F0-9]{40}"
	matched, err := regexp.MatchString(ethAddressRegex, operator)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidEthereumAddress
	}

	// todo : check existense of the node and read  the keystore from the file


}