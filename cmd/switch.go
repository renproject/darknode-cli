package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/renproject/darknode-cli/util"
	"github.com/renproject/phi"
	"github.com/urfave/cli"
)

// Commands for different actions to the darknode.
var (
	ActionStart   = "systemctl --user start darknode"
	ActionStop    = "systemctl --user stop darknode"
	ActionRestart = "systemctl --user restart darknode"
)

// updateServiceStatus can update status of the darknode service.
func updateServiceStatus(ctx *cli.Context, cmd string) error {
	tags := ctx.String("tags")
	name := ctx.Args().First()

	// Get the script we want to run depends on the command.
	var script, message string
	switch cmd {
	case "start":
		script, message = ActionStart, "started"
	case "stop":
		script, message = ActionStop, "stopped"
	case "restart":
		script, message = ActionRestart, "restarted"
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
		if errs[i] == nil {
			color.Green("[%v] has been %v", nodes[i], message)
		} else {
			color.Red("fail to %v [%v], err = %v", script, nodes[i])
		}
	})
	return util.HandleErrs(errs)
}
