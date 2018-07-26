package main

import (
	"os/exec"

	"github.com/republicprotocol/co-go"
	"github.com/urfave/cli"
)

func execScript(ctx *cli.Context) error {
	name := ctx.String("name")
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
		if len(nodes) == 0 {
			return ErrNoNodesFound
		}

		errs := make([]error, len(nodes))
		co.ForAll(nodes, func(i int) {
			errs[i] = execSingleNode(name, script)
		})

		for i := range errs {
			if errs[i] != nil {
				return err
			}
		}

		return nil
	}

	return ErrNameAndTags
}

func execSingleNode(name, script string) error {
	if script == "" {
		return ErrEmptyFilePath
	}
	nodeDirectory := Directory + "/darknodes/" + name
	keyPairPath := nodeDirectory + "/ssh_keypair"
	ip, err := getIp(nodeDirectory)
	if err != nil {
		return err
	}

	execCmd := exec.Command("ssh", "-i", keyPairPath, "ubuntu@"+ip, "-oStrictHostKeyChecking=no", "<", script)
	pipeToStd(execCmd)
	if err := execCmd.Start(); err != nil {
		return err
	}
	return execCmd.Wait()
}
