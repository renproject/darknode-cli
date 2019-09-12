package main

import (
	"errors"
	"io/ioutil"

	"github.com/republicprotocol/co-go"
	"github.com/urfave/cli"
)

// execScript execute a bash script on a darknode
// or a set of darknodes by the tags.
func execScript(ctx *cli.Context) error {
	name := ctx.Args().First()
	tags := ctx.String("tags")
	file := ctx.String("file")
	useSudo := ctx.Bool("sudo")
	script := ctx.String("script")

	if name == "" && tags == "" {
		cli.ShowCommandHelp(ctx, "exec")
		return ErrEmptyNameAndTags
	} else if name != "" && tags == "" {
		return execSingleNode(name, file, script, useSudo)
	} else if name == "" && tags != "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}

		errs := make([]error, len(nodes))
		co.ParForAll(nodes, func(i int) {
			errs[i] = execSingleNode(nodes[i], file, script, useSudo)
		})

		return handleErrs(errs)
	}

	return ErrNameAndTags
}

// execScript execute a bash script on a single darknode.
func execSingleNode(name, file, script string, useSudo bool) error {
	user := "darknode"
	if useSudo {
		user = "root"
	}
	if file != "" {
		script, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		return remoteUserRun(name, string(script), user)
	}

	if script != "" {
		return remoteUserRun(name, script, user)
	}

	return errors.New("please provide a script file or scripts to run ")
}
