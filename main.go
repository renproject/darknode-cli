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
	app.Name = "Darknode Deployer"
	app.Usage = "A command-line tool for managing Darknodes on Republic Protocol."
	app.Version = "1.0.0"

	upFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Value: "",
			Usage: "A unique name for your Darknode",
		},
		cli.StringFlag{
			Name:  "provider",
			Value: "",
			Usage: "The cloud service provider you want to use for your Darknode",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "",
			Usage: "The region you want to deploy to (default: random)",
		},
		cli.StringFlag{
			Name:  "instance",
			Value: "",
			Usage: "Instance type",
		},
		cli.StringFlag{
			Name:  "access-key",
			Value: "",
			Usage: "Access key for your AWS account, can be read from the default ~/.aws/credential file",
		},
		cli.StringFlag{
			Name:  "secret-key",
			Value: "",
			Usage: "Secret key for your AWS account, can be read from the default ~/.aws/credential file",
		},
		cli.StringFlag{
			Name:  "network",
			Value: "testnet",
			Usage: "The network you want to deploy your node to",
		},
	}

	destroyFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Value: "",
			Usage: "The name of the Darknode you want to destroy",
		},
		cli.BoolFlag{
			Name:  "skip",
			Usage: "Skip confirmation and begin destroying immediately",
		},
	}

	nameFlag := []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Value: "",
			Usage: "The name of the Darknode you want to operate",
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
			Flags: nameFlag,
			Action: func(c *cli.Context) error {
				return updateNode(c)
			},
		},
		{
			Name:  "ssh",
			Flags: nameFlag,
			Usage: "SSH into your Darknode",
			Action: func(c *cli.Context) error {
				return sshNode(c)
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
		fmt.Fprintf(c.App.Writer, "%scommand %q not found%s.\n", red, command, reset)
	}

	// Start the app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// updateNode update the Darknode to the latest release from master branch.
// This will restart the Darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.String("name")
	if name == "" {
		cli.ShowCommandHelp(ctx, "update")
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/" + name
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	updateScript := path.Join(os.Getenv("HOME"), ".darknode/scripts/update.sh")
	update, err := ioutil.ReadFile(updateScript)
	if err != nil {
		return err
	}
	keyPairPath := nodeDirectory + "/ssh_keypair"
	updateCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", string(update))
	pipeToStd(updateCmd)
	if err := updateCmd.Start(); err != nil {
		return err
	}

	if err := updateCmd.Wait(); err != nil {
		return err
	}
	fmt.Printf("%sDarknode has been updated to the latest version.%s", green, reset)

	return nil
}

// sshNode will ssh into the Darknode
func sshNode(ctx *cli.Context) error {
	name := ctx.String("name")
	if name == "" {
		cli.ShowCommandHelp(ctx, "ssh")
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/" + name
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	keyPairPath := nodeDirectory + "/ssh_keypair"
	ssh := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip)
	pipeToStd(ssh)
	if err := ssh.Start(); err != nil {
		return err
	}

	return ssh.Wait()
}

// listAllNodes will ssh into the Darknode
func listAllNodes() error {
	files, err := ioutil.ReadDir(Directory + "/darknodes")
	if err != nil {
		return err
	}

	fmt.Printf("%10s | %30s | %15s |\n", "name", "Address", "ip")

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

		fmt.Printf("%10s | %30s | %15s |\n", f.Name(), address, ip)
	}

	return nil
}
