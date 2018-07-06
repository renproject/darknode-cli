package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/urfave/cli"
)

// updateNode updates the Darknode to the latest release from master branch.
// This will restart the Darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.String("name")
	tag := ctx.String("tag")
	branch := ctx.String("branch")
	updateConfig := ctx.Bool("config")

	if name == "" && tag == "" {
		cli.ShowCommandHelp(ctx, "update")
		return ErrEmptyNodeName
	}

	// update a single darknode by its name
	if name != "" {
		if err := updateSingleNode(name, branch, updateConfig); err != nil {
			return err
		}
	}
	// Update a set of nodes by the tag
	if tag != "" {
		nodeNames, err := getNodesByTag(tag)
		if err != nil {
			return err
		}
		if len(nodeNames) == 0 {
			return ErrNoNodesFound
		}
		errs := make(chan error, len(nodeNames))
		dispatch.CoForAll(nodeNames, func(i int) {
			err := updateSingleNode(nodeNames[i], branch, updateConfig)
			if err != nil {
				errs <- err
			}
		})

		if len(errs) >= 1 {
			return <-errs
		}
	}

	return nil
}

func updateSingleNode(name, branch string, updateConfig bool) error {
	nodeDirectory := Directory + "/darknodes/" + name
	keyPairPath := nodeDirectory + "/ssh_keypair"
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}

	// Check if we need to update the node config
	if updateConfig {
		data, err := ioutil.ReadFile(nodeDirectory + "/config.json")
		if err != nil {
			return err
		}
		updateConfigScript := fmt.Sprintf(`echo "%s" >> $HOME/.darknode/config.json`, string(data))
		updateConfigCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", updateConfigScript)
		pipeToStd(updateConfigCmd)
		if err := updateConfigCmd.Start(); err != nil {
			return err
		}
		if err := updateConfigCmd.Wait(); err != nil {
			return err
		}
		fmt.Printf("%sConfig of [%s] has been updated to the local version.%s", GREEN, name, RESET)
	}

	updateScript := fmt.Sprintf(`
#!/usr/bin/env bash

cd ./go/src/github.com/republicprotocol/republic-go
sudo git stash
sudo git checkout %v
sudo git fetch origin %v
sudo git reset --hard origin/%v
cd cmd/darknode
go install
cd
sudo service darknode restart
`, branch, branch, branch)

	updateCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", updateScript)
	pipeToStd(updateCmd)
	if err := updateCmd.Start(); err != nil {
		return err
	}
	if err := updateCmd.Wait(); err != nil {
		return err
	}
	fmt.Printf("%s[%s] has been updated to the latest version on %s branch.%s", GREEN, name, branch, RESET)

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
