package main

import (
	"fmt"

	"github.com/renproject/phi"
	"github.com/republicprotocol/darknode-cli/util"
	"github.com/urfave/cli"
)

// Commands for different actions to the darknode.
var (
	ActionStart   = "systemctl --user start darknode"
	ActionStop    = "systemctl --user stop darknode"
	ActionRestart = "systemctl --user restart darknode"
)

// switchNode provide commands for basic operations to the darknode service.
func switchNode(ctx *cli.Context, cmd string) error {
	tags := ctx.String("tags")
	name := ctx.Args().First()

	// Get the script we want to run depends on the command.
	var script string
	switch cmd {
	case "start":
		script = ActionStart
	case "stop":
		script = ActionStop
	case "restart":
		script = ActionRestart
	default:
		panic(fmt.Sprintf("invalid switch command = %v", cmd))
	}

	// Parse the names of the node we want to operate
	nodes, err := util.ParseNodesFromNameAndTags(name, tags)
	if err != nil {
		return err
	}
	errs := make([]error, len(nodes))
	phi.ParForAll(nodes, func(i int) {
		errs[i] = util.RemoteRun(nodes[i], script)
	})
	return util.HandleErrs(errs)
}
