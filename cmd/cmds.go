package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/renproject/phi"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/urfave/cli"
)

var (
	// StartDarknode is command which will start the darknode service
	StartDarknode   = "systemctl --user start darknode"
	StopDarknode    = "systemctl --user stop darknode"
	RestartDarknode = "systemctl --user restart darknode"
)

// listAllNodes will list basic info of all the deployed darknodes.
// You can filter the results by the tags.
func listAllNodes(ctx *cli.Context) error {
	tags := ctx.String("tags")
	nodesNames, err := getNodesByTags(tags)
	if err != nil {
		return err
	}

	nodes := make([][]string, 0)
	for i := range nodesNames {
		tagFile := fmt.Sprintf("%v/darknodes/%v/tags.out", Directory, nodesNames[i])
		tags, err := ioutil.ReadFile(tagFile)
		if err != nil {
			continue
		}
		addressFile := fmt.Sprintf("%v/darknodes/%v/multiAddress.out", Directory, nodesNames[i])
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
		ethAddress, err := republicAddressToEthAddress(address)
		if err != nil {
			continue
		}
		prov, err := getProvider(nodesNames[i])
		if err != nil {
			continue
		}

		nodes = append(nodes, []string{nodesNames[i], address, ip, string(prov), string(tags), ethAddress.Hex()})
	}

	if len(nodes) > 0 {
		fmt.Printf("%-20s | %-30s | %-15s | %-10s | %-20s | %-45s \n", "name", "Address", "ip", "provider", "tags", "Ethereum Address")
		for i := range nodes {
			fmt.Printf("%-20s | %-30s | %-15s | %-10s | %-20s | %-45s\n", nodes[i][0], nodes[i][1], nodes[i][2], nodes[i][3], nodes[i][4], nodes[i][5])
		}
		return nil
	}

	return fmt.Errorf("%scannot find any node%s", RED, RESET)
}

// switchNode provide commands for basic operations to the darknode service.
func switchNode(ctx *cli.Context, cmd string) error {
	tags := ctx.String("tags")
	name := ctx.Args().First()

	// Get the function we want to run depends on the command.
	var f func(string) error
	switch cmd {
	case "start":
		f = startNode
	case "stop":
		f = stopNode
	case "restart":
		f = restartNode
	default:
		panic(fmt.Sprintf("invalid switch command = %v", cmd))
	}

	// Execute the operation on a single node or a set of nodes.
	if tags == "" && name == "" {
		return ErrEmptyNodeName
	} else if tags == "" && name != "" {
		return f(name)
	} else if tags != "" && name == "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}
		errs := make([]error, len(nodes))
		phi.ParForAll(nodes, func(i int) {
			errs[i] = startNode(nodes[i])
		})
		return handleErrs(errs)
	}

	return ErrNameAndTags
}

// startNode starts a single node.
func startNode(name string) error {
	err := remoteRun(name, StartDarknode)
	if err == nil {
		fmt.Printf("%s[%s] has been turned on.%s \n", GREEN, name, RESET)
	}
	return err
}

// stopNode stops a single node or a set of nodes by its tags.
func stopNode(name string) error {
	err := remoteRun(name, StopDarknode)
	if err == nil {
		fmt.Printf("%s[%s] has been turned off.%s \n", GREEN, name, RESET)
	}
	return err
}

// stopNode stops a single node or a set of nodes by its tags.
func restartNode(name string) error {
	err := remoteRun(name, RestartDarknode)
	if err == nil {
		fmt.Printf("%s[%s] has been restarted.%s \n", GREEN, name, RESET)
	}
	return err
}
