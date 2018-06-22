package main

import (
	"log"
	"os"
	"io/ioutil"
	"os/exec"
	"path"

	"github.com/urfave/cli"
)

var Directory =  path.Join(os.Getenv("HOME"), ".darknode")

func main() {
	// Create new cli application
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name : "name",
			Value : "",
			Usage : "name of your darknode so that you can easily distinguish between them",
		},
	}

	upFlags := []cli.Flag{
		cli.StringFlag{
			Name : "name",
			Value : "",
			Usage : "name of your darknode so that you can easily distinguish between them",
		},
		cli.StringFlag{
			Name:  "provider",
			Value: "AWS",
			Usage: "cloud service provider you want to use for your darknode, default to AWS",
		},
		cli.StringFlag{
			Name:  "region",
			Value: "",
			Usage: "region you want to deploy to",
		},
		cli.StringFlag{
			Name:  "instance",
			Value: "",
			Usage: "instance type",
		},
		cli.StringFlag{
			Name:  "access-key",
			Value: "",
			Usage: "access key for your AWS account",
		},
		cli.StringFlag{
			Name:  "secret-key",
			Value: "",
			Usage: "secret key for your AWS account",
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
			Usage: "secret key for your AWS account",
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
	updateScript  :=  path.Join(os.Getenv("HOME"), ".darknode/update.sh")
	update, err := ioutil.ReadFile(updateScript)
	if err != nil {
		return err
	}
	keyPairPath := path.Join(os.Getenv("HOME"), ".darknode/ssh_keypair")
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
	keyPairPath := path.Join(os.Getenv("HOME"), ".darknode/ssh_keypair")
	ssh := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip)
	pipeToStd(ssh)
	if err := ssh.Start(); err != nil {
		return err
	}

	return ssh.Wait()
}
// curl https://darknode.republicprotocol.com/darknode.sh -sSf | sh
