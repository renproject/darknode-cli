package main

import (
	"github.com/urfave/cli"
	"log"
	"os/exec"
)

var (
	green = "\x1B[1;32m"
	red   = "\x1B[1;31m"
	noc   = "\x1B[0m"
)

// destroyNode will tear down the deployed darknode, but keep the config file.
func destroyNode(ctx *cli.Context) error {
	// todo : how do we distinguish between AWS, digitalOcean and others.
	return destroyAwsNode()
}

// destroyAwsNode tear down the AWS instance.
func destroyAwsNode() error {
	log.Println("Destroying your darknode ...")
	init := exec.Command("./terraform", "destroy", "--force")
	if err := init.Run(); err != nil {
		return err
	}

	return init.Wait()
}
