package main

import (
	"errors"
	"io/ioutil"

	"github.com/renproject/phi"
	"github.com/republicprotocol/darknode-cli/util"
	"github.com/urfave/cli"
)

// execScript execute a bash script on a darknode or a set of darknodes by the tags.
func execScript(ctx *cli.Context) error {
	name := ctx.Args().First()
	tags := ctx.String("tags")
	file := ctx.String("file")
	script := ctx.String("script")

	// Parse the names of the node we want to operate
	nodes, err := util.ParseNodesFromNameAndTags(name, tags)
	if err != nil {
		return err
	}
	errs := make([]error, len(nodes))
	phi.ParForAll(nodes, func(i int) {
		errs[i] = execSingleNode(nodes[i], file, script)
	})
	return util.HandleErrs(errs)
}

// execScript execute a bash script on a single darknode.
func execSingleNode(name, file, script string) error {
	if file != "" {
		script, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		return util.RemoteRun(name, string(script))
	}

	if script != "" {
		return util.RemoteRun(name, script)
	}

	return errors.New("please provide a script file or scripts to run ")
}
