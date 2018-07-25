package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/urfave/cli"
)

// listAllNodes will ssh into the Darknode
func listAllNodes(ctx *cli.Context) error {
	tag := ctx.String("tag")

	files, err := ioutil.ReadDir(Directory + "/darknodes")
	if err != nil {
		return err
	}
	nodes := [][]string{}

	for _, f := range files {
		tagFile := Directory + "/darknodes/" + f.Name() + "/tags.out"
		tags, err := ioutil.ReadFile(tagFile)
		if err != nil {
			continue
		}
		if !strings.Contains(string(tags), tag) {
			continue
		}

		addressFile := Directory + "/darknodes/" + f.Name() + "/multiAddress.out"
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

		nodes = append(nodes, []string{f.Name(), address, ip, string(tags), ethAddress.Hex()})
	}

	if len(nodes) == 0 {
		return fmt.Errorf("%scannot find any node%s", RED, RESET)
	} else {
		fmt.Printf("%-20s | %-30s | %-15s | %-20s | %-45s \n", "name", "Address", "ip", "tags", "Ethereum Address")
		for i := range nodes {
			fmt.Printf("%-20s | %-30s | %-15s | %-20s | %-45s\n", nodes[i][0], nodes[i][1], nodes[i][2], nodes[i][3], nodes[i][4])
		}
	}

	return nil
}

// startNode starts a node by its name
func startNode(ctx *cli.Context) error {
	name := ctx.String("name")
	if name == "" {
		cli.ShowCommandHelp(ctx, "start")
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/" + name
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	startScript := "sudo systemctl start darknode"
	keyPairPath := nodeDirectory + "/ssh_keypair"
	startCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", startScript)
	pipeToStd(startCmd)
	if err := startCmd.Start(); err != nil {
		return err
	}
	if err := startCmd.Wait(); err != nil {
		return err
	}
	fmt.Printf("%s[%s] has been turned on.%s \n", GREEN, name, RESET)

	return nil
}

// stopNode stops a node by its name
func stopNode(ctx *cli.Context) error {

	name := ctx.String("name")
	if name == "" {
		cli.ShowCommandHelp(ctx, "stop")
		return ErrEmptyNodeName
	}
	nodeDirectory := Directory + "/darknodes/" + name
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	stopScript := "sudo systemctl stop darknode"
	keyPairPath := nodeDirectory + "/ssh_keypair"
	stopCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", stopScript)
	pipeToStd(stopCmd)
	if err := stopCmd.Start(); err != nil {
		return err
	}
	if err := stopCmd.Wait(); err != nil {
		return err
	}
	fmt.Printf("%s[%s] has been turned off.%s \n", GREEN, name, RESET)

	return nil
}
