package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/urfave/cli"
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

		nodes = append(nodes, []string{nodesNames[i], address, ip, string(tags), ethAddress.Hex()})
	}

	if len(nodes) > 0 {
		fmt.Printf("%-20s | %-30s | %-15s | %-20s | %-45s \n", "name", "Address", "ip", "tags", "Ethereum Address")
		for i := range nodes {
			fmt.Printf("%-20s | %-30s | %-15s | %-20s | %-45s\n", nodes[i][0], nodes[i][1], nodes[i][2], nodes[i][3], nodes[i][4])
		}
		return nil
	}

	return fmt.Errorf("%scannot find any node%s", RED, RESET)
}

// startNode starts a single node or a set of nodes by its tags.
func startNode(ctx *cli.Context) error {
	tags := ctx.String("tags")
	name := ctx.Args().First()

	if tags == "" && name == "" {
		return ErrEmptyNodeName
	} else if tags == "" && name != "" {
		return startSingleNode(name)
	} else if tags != "" && name == "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}
		errs := make([]error, len(nodes))
		co.ForAll(nodes, func(i int) {
			errs[i] = startSingleNode(nodes[i])
		})
		return handleErrs(errs)
	}

	return ErrNameAndTags
}

// startSingleNode starts a single node by its name
func startSingleNode(name string) error {
	nodePath := nodeDirPath(name)
	ip, err := getIp(nodePath)
	if err != nil {
		return err
	}
	startScript := "sudo systemctl start darknode"
	keyPairPath := nodePath + "/ssh_keypair"
	if err := run("ssh", "-i", keyPairPath, "darknode@"+ip, "-oStrictHostKeyChecking=no", startScript); err != nil {
		return err
	}
	fmt.Printf("%s[%s] has been turned on.%s \n", GREEN, name, RESET)

	return nil
}

// stopNode stops a single node or a set of nodes by its tags.
func stopNode(ctx *cli.Context) error {
	tags := ctx.String("tags")
	name := ctx.Args().First()

	if tags == "" && name == "" {
		return ErrEmptyNodeName
	} else if tags == "" && name != "" {
		return stopSingleNode(name)
	} else if tags != "" && name == "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}
		errs := make([]error, len(nodes))
		co.ForAll(nodes, func(i int) {
			errs[i] = stopSingleNode(nodes[i])
		})
		return handleErrs(errs)
	}

	return ErrNameAndTags
}

// stopSingleNode stops a single node by its name
func stopSingleNode(name string) error {
	if name == "" {
		return ErrEmptyNodeName
	}
	nodeDirectory := nodeDirPath(name)
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	stopScript := "sudo systemctl stop darknode"
	keyPairPath := nodeDirectory + "/ssh_keypair"
	if err := run("ssh", "-i", keyPairPath, "darknode@"+ip, "-oStrictHostKeyChecking=no", stopScript); err != nil {
		return err
	}

	fmt.Printf("%s[%s] has been turned off.%s \n", GREEN, name, RESET)
	return nil
}
