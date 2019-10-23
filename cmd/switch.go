package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/renproject/phi"
	"github.com/republicprotocol/darknode-cli/cmd/provider"
	"github.com/republicprotocol/darknode-cli/util"
	"github.com/republicprotocol/republic-go/identity"
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

// listAllNodes will list basic info of all the deployed darknodes.
// You can filter the results by the tags.
func listAllNodes(ctx *cli.Context) error {
	tags := ctx.String("tags")
	nodesNames, err := util.GetNodesByTags(tags)
	if err != nil {
		return err
	}

	nodes := make([][]string, 0)
	for i := range nodesNames {
		tagFile := filepath.Join(util.NodePath(nodesNames[i]), "tags.out")
		tags, err := ioutil.ReadFile(tagFile)
		if err != nil {
			continue
		}
		addressFile := filepath.Join(util.NodePath(nodesNames[i]), "multiAddress.out")
		data, err := ioutil.ReadFile(addressFile)
		if err != nil {
			continue
		}
		multi, err := identity.NewMultiAddressFromString(strings.TrimSpace(string(data)))
		if err != nil {
			continue
		}
		address, err := multi.ValueForProtocol(identity.RepublicCode)
		if err != nil {
			continue
		}
		ip, err := multi.ValueForProtocol(identity.IP4Code)
		if err != nil {
			continue
		}
		// FIXME : GET THE CORRECT ETHEREUM ADDRESS
		ethAddress := "0xThisIsAnEthereumAddress"
		prov, err := provider.GetProvider(nodesNames[i])
		if err != nil {
			continue
		}

		nodes = append(nodes, []string{nodesNames[i], address, ip, string(prov), string(tags), ethAddress})
	}

	if len(nodes) > 0 {
		fmt.Printf("%-20s | %-30s | %-15s | %-10s | %-20s | %-45s \n", "name", "Address", "ip", "provider", "tags", "Ethereum Address")
		for i := range nodes {
			fmt.Printf("%-20s | %-30s | %-15s | %-10s | %-20s | %-45s\n", nodes[i][0], nodes[i][1], nodes[i][2], nodes[i][3], nodes[i][4], nodes[i][5])
		}
		return nil
	}

	return util.RedError("cannot find any node")
}
