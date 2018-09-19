package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/urfave/cli"
)

// updateNode updates the Darknode to the latest release from master branch.
// This will restart the Darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	updateConfig := ctx.Bool("config")
	tags := ctx.String("tags")
	branch := ctx.String("branch")

	if tags == "" && name == "" {
		return ErrEmptyNodeName
	} else if tags == "" && name != "" {
		return updateSingleNode(name, branch, updateConfig)
	} else if tags != "" && name == "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}
		errs := make([]error, len(nodes))
		co.ForAll(nodes, func(i int) {
			errs[i] = updateSingleNode(nodes[i], branch, updateConfig)
		})
		return handleErrs(errs)
	}

	return ErrNameAndTags
}

func updateSingleNode(name, branch string, updateConfig bool) error {
	nodeDir := nodeDirectory(name)
	keyPairPath := nodeDir + "/ssh_keypair"
	ip, err := getIp(nodeDir)
	if err != nil {
		return err
	}

	// Check if we need to update the node config
	if updateConfig {
		data, err := ioutil.ReadFile(nodeDir + "/config.json")
		if err != nil {
			return err
		}
		updateConfigScript := fmt.Sprintf(`echo '%s' > $HOME/.darknode/config.json`, string(data))
		updateConfigCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", updateConfigScript)
		pipeToStd(updateConfigCmd)
		if err := updateConfigCmd.Start(); err != nil {
			return err
		}
		if err := updateConfigCmd.Wait(); err != nil {
			return err
		}
		fmt.Printf("%sConfig of [%s] has been updated to the local version.%s\n", GREEN, name, RESET)
	}

	// Default branch is depends on the network parameter.
	if branch == "" {
		config, err := config.NewConfigFromJSONFile(nodeDir + "/config.json")
		if err != nil {
			return err
		}
		switch config.Ethereum.Network {
		case "mainnet":
			branch = "master"
		case "testnet":
			branch = "develop"
		case "falcon", "nightly":
			branch = "nightly"
		default:
			panic("unknown network")
		}
	}

	updateScript := fmt.Sprintf(`
#!/usr/bin/env bash

cd ./go/src/github.com/republicprotocol/republic-go
sudo git reset --hard HEAD
sudo git clean -f -d
sudo git checkout %v
sudo git fetch origin %v
sudo git reset --hard origin/%v
cd cmd/darknode
go install
cd
sudo service darknode restart

curl -s 'https://darknode.republicprotocol.com/auto-updater.sh' > .darknode/updater.sh
sudo service darknode-updater restart
`, branch, branch, branch)

	updateCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", updateScript)
	pipeToStd(updateCmd)
	if err := updateCmd.Start(); err != nil {
		return err
	}
	if err := updateCmd.Wait(); err != nil {
		return err
	}

	fmt.Printf("%s[%s] has been updated to the latest version on %s branch.%s \n", GREEN, name, branch, RESET)
	return nil
}

// sshNode will ssh into the Darknode
func sshNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	if name == "" {
		cli.ShowCommandHelp(ctx, "ssh")
		return ErrEmptyNodeName
	}
	nodeDirectory := nodeDirectory(name)
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
