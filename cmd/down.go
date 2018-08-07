package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/republicprotocol/co-go"
	"github.com/urfave/cli"
)

// destroyNode tears down the deployed darknode(s).
func destroyNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	force := ctx.Bool("force")
	tags := ctx.String("tags")

	if tags == "" && name == "" {
		cli.ShowCommandHelp(ctx, "down")
		return ErrEmptyNodeName
	} else if tags == "" && name != "" {
		return destroySingleNode(name, force)
	} else if tags != "" && name == "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}
		errs := make([]error, len(nodes))
		co.ForAll(nodes, func(i int) {
			errs[i] = destroySingleNode(nodes[i], force)
		})
		return handleErrs(errs)
	}

	return ErrNameAndTags
}

// destroySingleNode tears down a single darknode by its name.
func destroySingleNode(name string, force bool) error {
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
	cmd := fmt.Sprintf("cd %v && terraform destroy --force && rm -rf %v", nodeDirectory, nodeDirectory)
	destroy := exec.Command("bash", "-c", cmd)
	pipeToStd(destroy)
	if err := destroy.Start(); err != nil {
		return err
	}

	return destroy.Wait()
}
