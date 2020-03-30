package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/renproject/darknode-cli/util"
	"github.com/renproject/phi"
	"github.com/urfave/cli"
)

// updateNode updates the Darknode to the latest release or the version specified
// by user.
func updateNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	tags := ctx.String("tags")
	version := strings.TrimSpace(ctx.String("version"))
	nodes, err := util.ParseNodesFromNameAndTags(name, tags)
	if err != nil {
		return err
	}

	// Use latest version if user doesn't provide a version number
	if version == "" {
		version, err = util.LatestReleaseVersion()
		if err != nil {
			return err
		}
	}
	// Check if the target release exists on github
	color.Green("Verifying darknode release ...")
	if err := validateVersion(version); err != nil {
		return err
	}

	color.Green("Updating darknodes...")
	errs := make([]error, len(nodes))
	phi.ParForAll(nodes, func(i int) {
		errs[i] = updateSingleNode(nodes[i], version)
	})
	return util.HandleErrs(errs)
}

func updateSingleNode(name, ver string) error {
	v, _ := util.Version(name)
	curVersion, err := version.NewVersion(strings.TrimSpace(v))
	if err != nil {
		return err
	}
	newVersion, err := version.NewVersion(strings.TrimSpace(ver))
	res := curVersion.Compare(newVersion)
	switch res {
	case 0:
		color.Green("darknode [%v] is running version [%v] already.", name, ver)
	case 1:
		color.Red("darknode [%v] is running with version %v, you cannot downgrade to a lower version %v", name, curVersion.String(), newVersion.String())
	default:
		url := fmt.Sprintf("https://github.com/renproject/darknode-release/releases/download/%v", ver)
		script := fmt.Sprintf(`mv ~/.darknode/bin/darknode ~/.darknode/bin/darknode-backup && 
curl -sL %v/darknode > ~/.darknode/bin/darknode && 
curl -sL %v/migration > ~/.darknode/bin/migration &&
chmod +x ~/.darknode/bin/darknode && 
chmod +x ~/.darknode/bin/migration &&
systemctl --user stop darknode &&
cp -a ~/.darknode/db/. ~/.darknode/db_bak/ &&
~/.darknode/bin/migration &&
rm -rf ~/.darknode/db &&
mv ~/.darknode/db_bak ~/.darknode/db &&
echo %v > ~/.darknode/version &&
systemctl --user restart darknode`, url, url, ver)
		err = util.RemoteRun(name, script)
		if err != nil {
			color.Red("cannot update darknode %v, error = %v", name, err)
		} else {
			color.Green("[%s] has been updated to version %v", name, ver)
		}
	}
	return nil
}

func validateVersion(version string) error {
	url := fmt.Sprintf("https://api.github.com/repos/renproject/darknode-release/releases/tags/%v", version)
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	// Check the status code of the response
	switch response.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("cannot find release [%v] on github", version)
	default:
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("cannot connect to github, code= %v, err = %v", response.StatusCode, string(data))
	}
}
