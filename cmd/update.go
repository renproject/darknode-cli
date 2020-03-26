package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
	"github.com/renproject/darknode-cli/util"
	"github.com/renproject/phi"
	"github.com/urfave/cli"
)

// updateNode updates the Darknode to the latest release. It can also be used
// to update the config file of the darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	if err := util.ValidateNodeName(name); err != nil {
		return err
	}
	tags := ctx.String("tags")
	version := ctx.String("version")

	nodes, err := util.ParseNodesFromNameAndTags(name, tags)
	if err != nil {
		return err
	}
	errs := make([]error, len(nodes))
	phi.ParForAll(nodes, func(i int) {
		errs[i] = updateSingleNode(nodes[i], version)
	})
	return util.HandleErrs(errs)
}

func updateSingleNode(name, ver string) error {
	url := "https://www.github.com/renproject/darknode-release/releases/latest/download/darknode"
	if ver != "" {
		if err := validateVersion(ver); err != nil {
			return err
		}
		url = fmt.Sprintf("https://github.com/renproject/darknode-release/releases/download/%v/darknode", ver)
	}

	script := fmt.Sprintf(`rm ~/.darknode/bin/darknode && curl -sL %v > ~/.darknode/bin/darknode && chmod +x ~/.darknode/bin/darknode && systemctl --user restart darknode`, url)
	if err := util.RemoteRun(name, script); err != nil {
		return err
	}
	if ver == "" {
		color.Green("[%s] has been updated to the latest version", name)
	} else {
		color.Green("[%s] has been updated to version %v", name, ver)
	}
	return nil
}

func validateVersion(version string) error {
	url := fmt.Sprintf("https://api.github.com/repos/renproject/darknode-release/releases/tags/%v", version)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("cannot find release [%v] on github", version)
	}
	if response.StatusCode == http.StatusOK {
		return nil
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("cannot connect to github, err = %v", string(data))
}
