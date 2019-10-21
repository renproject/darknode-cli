package main

import (
	"os"
	"path/filepath"

	"github.com/republicprotocol/co-go"
	"github.com/urfave/cli"
)

// execScript execute a bash script on a darknode
// or a set of darknodes by the tags.
func execScript(ctx *cli.Context) error {
	name := ctx.Args().First()
	tags := ctx.String("tags")
	script := ctx.String("script")

	if name == "" && tags == "" {
		cli.ShowCommandHelp(ctx, "update")
		return ErrEmptyNameAndTags
	} else if name != "" && tags == "" {
		return execSingleNode(name, script)
	} else if name == "" && tags != "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}

		errs := make([]error, len(nodes))
		co.ForAll(nodes, func(i int) {
			errs[i] = execSingleNode(name, script)
		})

		return handleErrs(errs)
	}

	return ErrNameAndTags
}

// execScript execute a bash script on a single darknode.
func execSingleNode(name, script string) error {
	if script == "" {
		return ErrEmptyFilePath
	}
	nodeDirectory := nodePath(name)
	keyPairPath := nodeDirectory + "/ssh_keypair"
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	filePath := filepath.Join(cwd, script)

	// todo : why this not working?
	return run("ssh", "-i", keyPairPath, "darknode@"+ip, "'bash -s'", "", filePath)
}
