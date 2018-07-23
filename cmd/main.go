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
		Usage: "A unique human-readable `string` for identifying the Darknode",
	}
	tagFlag := cli.StringFlag{
		Name:  "tag",
		Usage: "A human-readable `string` for identifying groups of Darknodes",
	}
	tagsFlag := cli.StringFlag{
		Name:  "tags",
		Usage: "Multiple human-readable comma separated `strings` for identifying groups of Darknodes",
	}

	// Flag for each command
	upFlags := []cli.Flag{
		nameFlag, tagsFlag,
		cli.StringFlag{
			Name:  "keystore",
			Usage: "An optional keystore `file` that will be used for the Darknode",
		},
		cli.StringFlag{
			Name:  "passphrase",
			Usage: "An optional `secret` for decrypting the keystore file",
		},
		cli.StringFlag{
			Name:  "config",
			Usage: "An optional configuration `file` for the Darknode",
		},
		cli.StringFlag{
			Name:  "network",
			Value: "testnet",
			Usage: "Darkpool network of your node",
		},

		// AWS flags
		cli.BoolFlag{
			Name:  "aws",
			Usage: "AWS will be used to provision the Darknode",
		},
		cli.StringFlag{
			Name:  "aws-access-key",
			Usage: "AWS access `key` for programmatic access",
		},
		cli.StringFlag{
			Name:  "aws-secret-key",
			Usage: "AWS secret `key` for programmatic access",
		},
		cli.StringFlag{
			Name:  "aws-region",
			Usage: "An optional AWS region (default: random)",
		},
		cli.StringFlag{
			Name:  "aws-instance",
			Value: "t2.medium",
			Usage: "An optional AWS EC2 instance type",
		},
		cli.StringFlag{
			Name:  "aws-elastic-ip",
			Usage: "An optional allocation ID for an elastic IP address",
		},

		// Digital Ocean flags
		cli.BoolFlag{
			Name:  "digitalocean",
			Usage: "Digital Ocean will be used to provision the Darknode",
		},
	}

	updateFlags := []cli.Flag{
		nameFlag, tagFlag,
		cli.StringFlag{
			Name:  "branch, b",
			Value: "master",
			Usage: "Release `branch` used to update the software",
		},
		cli.BoolFlag{
			Name:  "config, c",
			Usage: "An optional configuration `file` used to update the configuration",
		},
	}

	destroyFlags := []cli.Flag{
		nameFlag,
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Force destruction without interactive prompts",
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
			Usage:   "Destroy one of your Darknode",
			Aliases: []string{"down"},
			Flags:   destroyFlags,
			Action: func(c *cli.Context) error {
				return destroyNode(c)
			},
		},
		{
			Name:  "update",
			Usage: "Update your Darknodes to the latest software and configuration",
			Flags: updateFlags,
			Action: func(c *cli.Context) error {
				return updateNode(c)
			},
		},
		{
			Name:  "ssh",
			Flags: []cli.Flag{nameFlag},
			Usage: "SSH into one of your Darknode",
			Action: func(c *cli.Context) error {
				return sshNode(c)
			},
		},
		{
			Name:  "start",
			Flags: []cli.Flag{nameFlag},
			Usage: "Start one of your Darknodes from a suspended state",
			Action: func(c *cli.Context) error {
				return startNode(c)
			},
		},
		{
			Name:  "stop",
			Flags: []cli.Flag{nameFlag},
			Usage: "Stop one of your Darknodes by putting it into a suspended state",
			Action: func(c *cli.Context) error {
				return stopNode(c)
			},
		},
		{
			Name:  "list",
			Usage: "List all of your Darknodes",
			Flags: []cli.Flag{tagFlag},
			Action: func(c *cli.Context) error {
				return listAllNodes(c)
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
func listAllNodes(ctx *cli.Context) error {
	tag := ctx.String("tag")

	files, err := ioutil.ReadDir(Directory + "/darknodes")
	if err != nil {
		return err
	}
	nodes := [][]string{}

	for _, f := range files {
		tagFile := Directory + "/darknodes/" + f.Name() + "/tags.out"
		tags, err := ioutil.ReadFile(tagFile)
		if err != nil {
			continue
		}
		if !strings.Contains(string(tags), tag) {
			continue
		}

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
		ethAddress, err := republicAddressToEthAddress(address)
		if err != nil {
			continue
		}

		nodes = append(nodes, []string{f.Name(), address, ip, string(tags) , ethAddress.Hex()})
	}

	if len(nodes) == 0 {
		return fmt.Errorf("%scannot find any node%s", RED, RESET)
	} else {
		fmt.Printf("%-20s | %-30s | %-15s | %-20s | %-45s \n", "name", "Address", "ip", "tags", "Ethereum Address")
		for i := range nodes {
			fmt.Printf("%-20s | %-30s | %-15s | %-20s | %-45s\n", nodes[i][0], nodes[i][1], nodes[i][2], nodes[i][3],  nodes[i][4])
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
	fmt.Printf("%s[%s] has been turned on.%s \n", GREEN, name, RESET)

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
	fmt.Printf("%s[%s] has been turned off.%s \n", GREEN, name, RESET)

	return nil
}
