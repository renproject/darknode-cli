package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/fatih/color"
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
		script := fmt.Sprintf("echo '%s' > $HOME/.darknode/config.json", string(data))
		if err := util.RemoteRun(name, script); err != nil {
			return err
		}
		color.Green("Config of [%s] has been updated to the local version.", name)
	}
	updateScript := "$HOME/.darknode/bin/update.sh"
	if err := util.RemoteRun(name, updateScript); err != nil {
		return err
	}
	color.Green("[%s] has been updated to the latest version.", name)
	return nil
}
