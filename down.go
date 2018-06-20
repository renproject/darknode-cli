package main

import (
	"github.com/urfave/cli"
	"log"
	"os/exec"
)

// destroyNode will tear down the deployed darknode, but keep the config file.
func destroyNode(ctx *cli.Context) error {
	// FIXME : currently it only supports tear down AWS deployment.
	// Needs to figure out way which suits for all kinds of cloud service.
	return destroyAwsNode()
}

// destroyAwsNode tear down the AWS instance.
func destroyAwsNode() error {
	log.Println("Destroying your darknode ...")
	destroy := exec.Command("./terraform", "destroy", "--force")
	pipeToStd(destroy)
	if err := destroy.Start(); err != nil {
		return err
	}

	return destroy.Wait()
}
