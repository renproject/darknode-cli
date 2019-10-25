package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/renproject/phi"
	"github.com/republicprotocol/darknode-cli/cmd/provider"
	"github.com/republicprotocol/darknode-cli/util"
	"github.com/urfave/cli"
)

// listAllNodes will list information of deployed darknodes. Results can be filtered by the tags.
func listAllNodes(ctx *cli.Context) error {
	tags := ctx.String("tags")
	nodesNames, err := util.GetNodesByTags(tags)
	if err != nil {
		return err
	}

	nodes := make([][]string, len(nodesNames))
	errs := map[string]error{}
	phi.ParForAll(nodesNames, func(i int) {
		name := nodesNames[i]
		var err error
		nodes[i], err = func() ([]string, error) {
			id, err := util.ID(name)
			if err != nil {
				return nil, err
			}
			ip, err := util.IP(name)
			if err != nil {
				return nil, err
			}
			provider, err := provider.GetProvider(name)
			if err != nil {
				return nil, err
			}
			tagFile := filepath.Join(util.NodePath(nodesNames[i]), "tags.out")
			tags, err := ioutil.ReadFile(tagFile)
			if err != nil {
				return nil, err
			}
			ethAddr, err := id.ToEthereumAddress()
			if err != nil {
				return nil, err
			}
			return []string{name, id.String(), ip, provider, string(tags), ethAddr.Hex()}, nil
		}()
		if err != nil {
			errs[name] = err
		}
	})
	if len(errs) == len(nodesNames) {
		return fmt.Errorf("cannot find any node")
	}

	fmt.Printf("%-20s | %-30s | %-15s | %-8s | %-15s | %-45s \n", "name", "id", "ip", "provider", "tags", "ethereum address")
	for _, node := range nodes {
		fmt.Printf("%-20s | %-30s | %-15s | %-8s | %-15s | %-45s\n", node[0], node[1], node[2], node[3], node[4], node[5])
	}
	return nil
}
