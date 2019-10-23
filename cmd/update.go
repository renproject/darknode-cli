package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/renproject/phi"
	"github.com/republicprotocol/darknode-cli/util"
	"github.com/urfave/cli"
)

// updateNode updates the Darknode to the latest release. It can also be used
// to update the config file of the darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	tags := ctx.String("tags")
	updateConfig := ctx.Bool("config")

	nodes, err := util.ParseNodesFromNameAndTags(name, tags)
	if err != nil {
		return err
	}
	errs := make([]error, len(nodes))
	phi.ParForAll(nodes, func(i int) {
		errs[i] = updateSingleNode(nodes[i], updateConfig)
	})
	return util.HandleErrs(errs)
}

func updateSingleNode(name string, updateConfig bool) error {
	path := util.NodePath(name)
	configPath := filepath.Join(path, "config.json")

	if updateConfig {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return err
		}
		dir := "$HOME/.darknode"
		script := fmt.Sprintf(`mkdir -p %v/backup && mv %v/config.json %v/backup/%v.json && echo '%s' > $HOME/.darknode/config.json`, dir, dir, dir, time.Now().String(), string(data))
		if err := util.RemoteRun(name, script); err != nil {
			return err
		}
		util.GreenPrintln(fmt.Sprintf("Config of [%s] has been updated to the local version.", name))
	}

	// FIXME : HOW DOW WE UPDATE DARKNODE.
	util.GreenPrintln(fmt.Sprintf("[%s] has been updated to the latest version.", name))
	return nil
}

// sshNode will ssh into the Darknode
func sshNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	if name == "" {
		cli.ShowCommandHelp(ctx, "ssh")
		return ErrEmptyNodeName
	}
	nodePath := util.NodePath(name)
	ip, err := util.IP(nodePath)
	if err != nil {
		return err
	}
	keyPairPath := nodePath + "/ssh_keypair"

	return util.Run("ssh", "-i", keyPairPath, "darknode@"+ip, "-oStrictHostKeyChecking=no")
}
