package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/urfave/cli"
)

// Directory is the directory address of the deployer and Darknodes.
var Directory = path.Join(os.Getenv("HOME"), ".darknode")

func main() {
	// Create new cli application
	app := cli.NewApp()
	app.Name = "Darknode CLI"
	app.Usage = "A command-line tool for managing Darknodes."
	app.Version = "1.1.0"

	// Define some of the flags
	nameFlag := cli.StringFlag{
		Name:  "name",
		Value: "",
		Usage: "Unique name of the Darknode",
	}
	tagFlag := cli.StringFlag{
		Name:  "tag",
		Usage: "Tag of darknodes ",
	}
	tagsFlag := cli.StringFlag{
		Name:  "tags",
		Usage: "Tags for the Darknode ",
	}

	// Flag for each command
	upFlags := []cli.Flag{
		nameFlag, tagsFlag,
		cli.StringFlag{
			Name:  "keystore",
			Usage: "Name of the keystore `file` for the Darknode",
		},
		cli.StringFlag{
			Name:  "passphrase",
			Usage: "Passphrase for decrypting the keystore file",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "Name of the configuration `file` for the Darknode",
		},
		cli.BoolFlag{
			Name:  "aws",
			Usage: "Use AWS to provision the Darknode",
		},
		cli.BoolFlag{
			Name:  "digitalocean",
			Usage: "Use digital-ocean to provision the Darknode",
		},
		cli.StringFlag{
			Name:  "aws-region",
			Usage: "AWS region for the Darknode ",
		},
		cli.StringFlag{
			Name:  "aws-instance",
			Value: "t2.small",
			Usage: "AWS EC2 instance type for the Darknode",
		},
		cli.StringFlag{
			Name:  "aws-access-key",
			Usage: "AWS access `key` (defaults to $HOME/.aws/credential)",
		},
		cli.StringFlag{
			Name:  "aws-secret-key",
			Usage: "AWS secret `key` (defaults to $HOME/.aws/credential)",
		},
		cli.StringFlag{
			Name:  "aws-allocation-id",
			Usage: "Allocation ID of the elastic IP you want to associate",
		},
	}

	updateFlags := []cli.Flag{
		nameFlag, tagFlag,
		cli.BoolFlag{
			Name:  "config, c",
			Usage: "update the node config to the local one",
		},
		cli.StringFlag{
			Name:  "branch, b",
			Value: "master",
			Usage: "branch name of republic-go repo",
		},
	}

	destroyFlags := []cli.Flag{
		nameFlag,
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Force the Darknode to be destroyed without interactive prompts",
		},
	}

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "Deploy a new Darknode",
			Flags: upFlags,
			Action: func(c *cli.Context) error {
				return deployNode(c)
			},
		},
		{
			Name:    "destroy",
			Usage:   "Tear down the Darknode and clean-up files",
			Aliases: []string{"down"},
			Flags:   destroyFlags,
			Action: func(c *cli.Context) error {
				return destroyNode(c)
			},
		},
		{
			Name:  "update",
			Usage: "Update your Darknode to the latest release",
			Flags: updateFlags,
			Action: func(c *cli.Context) error {
				return updateNode(c)
			},
		},
		{
			Name:  "ssh",
			Flags: []cli.Flag{nameFlag},
			Usage: "SSH into your Darknode",
			Action: func(c *cli.Context) error {
				return sshNode(c)
			},
		},
		{
			Name:  "start",
			Flags: []cli.Flag{nameFlag},
			Usage: "Start a Darknode by its name",
			Action: func(c *cli.Context) error {
				return startNode(c)
			},
		},
		{
			Name:  "stop",
			Flags: []cli.Flag{nameFlag},
			Usage: "Stop a Darknode by its name",
			Action: func(c *cli.Context) error {
				return stopNode(c)
			},
		},
		{
			Name:  "list",
			Usage: "List all your deployed Darknodes",
			Action: func(c *cli.Context) error {
				return listAllNodes()
			},
		},
	}

	// Show error message and display the help page for the app
	app.CommandNotFound = func(c *cli.Context, command string) {
		cli.ShowAppHelp(c)
		fmt.Fprintf(c.App.Writer, "%scommand %q not found%s.\n", RED, command, RESET)
	}

	// Start the app
	err := app.Run(os.Args)
	if err != nil {
		// Remove the timestamp for error message
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
		log.Fatal(err)
	}
}

// listAllNodes will ssh into the Darknode
func listAllNodes() error {
	files, err := ioutil.ReadDir(Directory + "/darknodes")
	if err != nil {
		return err
	}
	nodes := [][]string{}

	for _, f := range files {
		addressFile := Directory + "/darknodes/" + f.Name() + "/multiAddress.out"
		data, err := ioutil.ReadFile(addressFile)
		if err != nil {
			continue
		}
		multi, err := identity.NewMultiAddressFromString(strings.TrimSpace(string(data)))
		if err != nil {
			continue
		}
		address, err := multi.ValueForProtocol(identity.RepublicCode)
		if err != nil {
			continue
		}
		ip, err := multi.ValueForProtocol(identity.IP4Code)
		if err != nil {
			continue
		}

		tagFile := Directory + "/darknodes/" + f.Name() + "/tags.out"
		tags, err := ioutil.ReadFile(tagFile)
		if err != nil {
			continue
		}

		nodes = append(nodes, []string{f.Name(), address, ip, string(tags)})
	}

	if len(nodes) == 0 {
		return fmt.Errorf("%scannot find any node%s", RED, RESET)
	} else {
		fmt.Printf("%-20s | %-30s | %-15s | %-20s \n", "name", "Address", "ip", "tags")
		for i := range nodes {
			fmt.Printf("%-20s | %-30s | %-15s | %-20s \n", nodes[i][0], nodes[i][1], nodes[i][2], nodes[i][3])
		}
	}

	return nil
}

// startNode starts a node by its name
func startNode(ctx *cli.Context) error {
	name := ctx.String("name")
	if name == "" {
		cli.ShowCommandHelp(ctx, "start")
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/" + name
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	startScript := "sudo systemctl start darknode"
	keyPairPath := nodeDirectory + "/ssh_keypair"
	startCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", startScript)
	pipeToStd(startCmd)
	if err := startCmd.Start(); err != nil {
		return err
	}
	if err := startCmd.Wait(); err != nil {
		return err
	}
	fmt.Printf("%sDarknode has been turned on.%s \n", GREEN, RESET)

	return nil
}

// stopNode stops a node by its name
func stopNode(ctx *cli.Context) error {

	name := ctx.String("name")
	if name == "" {
		cli.ShowCommandHelp(ctx, "stop")
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/" + name
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	stopScript := "sudo systemctl stop darknode"
	keyPairPath := nodeDirectory + "/ssh_keypair"
	stopCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", stopScript)
	pipeToStd(stopCmd)
	if err := stopCmd.Start(); err != nil {
		return err
	}
	if err := stopCmd.Wait(); err != nil {
		return err
	}
	fmt.Printf("%sDarknode has been turned off.%s \n", GREEN, RESET)

	return nil
}
