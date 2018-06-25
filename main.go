package main

import (
	"log"
	"os"
	"io/ioutil"
	"os/exec"
	"path"
	"fmt"
	"strings"

	"github.com/urfave/cli"
	"github.com/republicprotocol/republic-go/identity"
)

var Directory =  path.Join(os.Getenv("HOME"), ".darknode")

func main() {
	// Create new cli application
	app := cli.NewApp()

	upFlags := []cli.Flag{
		cli.StringFlag{
			Name : "name",
			Value : "",
			Usage : "name of your darknode so that you can easily distinguish between them",
		},
		cli.StringFlag{
			Name:  "provider",
			Value: "AWS",
			Usage: "cloud service provider you want to use for your darknode. (default to AWS)",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "",
			Usage: "region you want to deploy to. (default to random)",
		},
		cli.StringFlag{
			Name:  "instance",
			Value: "",
			Usage: "instance type. (default to `t2.small`)",
		},
		cli.StringFlag{
			Name:  "access-key",
			Value: "",
			Usage: "access key for your AWS account, can be read from the default .aws/credential file",
		},
		cli.StringFlag{
			Name:  "secret-key",
			Value: "",
			Usage: "secret key for your AWS account, can be read from the default .aws/credential file",
		},
		cli.StringFlag{
			Name:  "network",
			Value: "testnet",
			Usage: "which network you want to deploy your node to.(default to `testnet`)",
		},
	}

	destroyFlags := []cli.Flag{
		cli.StringFlag{
			Name : "name",
			Value : "",
			Usage : "name of your darknode so that you can easily distinguish between them",
		},
		cli.BoolFlag{
			Name:  "skip",
			Usage: "secret key for your AWS account.(default to `false`)",
		},
	}

	nameFlag := []cli.Flag{
		cli.StringFlag{
			Name : "name",
			Value : "",
			Usage : "name of your darknode so that you can easily distinguish between them",
		},
	}

	// Define sub-commands
	app.Commands = []cli.Command{
		{
			Name:  "up",
			Usage: "deploying a new darknode",
			Flags: upFlags,
			Action: func(c *cli.Context) error {
				return deployNode(c)
			},
		},
		{
			Name:  "destroy",
			Usage: "tear down the darknode and clean up everything",
			Aliases: []string{"down"},
			Flags: destroyFlags,
			Action: func(c *cli.Context) error {
				return destroyNode(c)
			},
		},
		{
			Name:  "update",
			Usage: "update your darknode to the latest release",
			Flags: nameFlag,
			Action: func(c *cli.Context) error {
				return updateNode(c)
			},
		},
		{
			Name:  "ssh",
			Flags: nameFlag,
			Usage: "ssh into your cloud service instance",
			Action: func(c *cli.Context) error {
				return sshNode(c)
			},
		},
		{
			Name:  "list",
			Usage: "list all the d",
			Action: func(c *cli.Context) error {
				return listAllNodes(c)
			},
		},
	}

	// Start the app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// updateNode update the darknode to the latest release from master branch.
// This will restart the darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.String("name")
	if name == ""{
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/"+  name
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	updateScript  :=  path.Join(os.Getenv("HOME"), ".darknode/scripts/update.sh")
	update, err := ioutil.ReadFile(updateScript)
	if err != nil {
		return err
	}
	keyPairPath := nodeDirectory + "/ssh_keypair"
	updateCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no",  string(update))
	pipeToStd(updateCmd)
	if err := updateCmd.Start(); err != nil {
		return err
	}

	return updateCmd.Wait()
}

// sshNode will ssh into the darknode
func sshNode(ctx *cli.Context) error {
	name := ctx.String("name")
	if name == ""{
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/"+  name
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

// listAllNodes will ssh into the darknode
func listAllNodes(ctx *cli.Context) error {
	files, err := ioutil.ReadDir(Directory + "/darknodes")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%10s | %30s | %15s |\n", "name","Address", "ip" )

	for _, f := range files {
		addressFile :=  Directory + "/darknodes/" + f.Name() + "/multiAddress.out"
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

		fmt.Printf( "%10s | %30s | %15s |\n", f.Name(), address, ip)
	}

	return nil
}
