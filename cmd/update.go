package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/google/go-github/v31/github"
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
	force := ctx.Bool("downgrade")
	version := strings.TrimSpace(ctx.String("version"))
	nodes, err := util.ParseNodesFromNameAndTags(name, tags)
	if err != nil {
		return err
	}

	// Use latest version if user doesn't provide a version number
	if version == "" {
		version, err = util.LatestStableRelease()
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
		errs[i] = updateSingleNode(nodes[i], version, force)
	})
	return util.HandleErrs(errs)
}

func updateSingleNode(name, ver string, force bool) error {
	v := util.Version(name)
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
		if !force {
			color.Red("darknode [%v] is running with version %v, you cannot downgrade to a lower version %v", name, curVersion.String(), newVersion.String())
			return nil
		}
		if err := update(name, ver); err != nil {
			color.Red("cannot downgrade darknode %v, error = %v", name, err)
		} else {
			color.Green("[%s] has been downgraded to version %v", name, ver)
		}
	default:
		if err := update(name, ver); err != nil {
			color.Red("cannot update darknode %v, error = %v", name, err)
		} else {
			color.Green("[%s] has been updated to version %v", name, ver)
		}
	}
	return nil
}

func update(name, ver string) error {
	url := fmt.Sprintf("https://github.com/renproject/darknode-release/releases/download/%v", ver)
	script := fmt.Sprintf(`mv ~/.darknode/bin/darknode ~/.darknode/bin/darknode-backup && 
curl -sL %v/darknode > ~/.darknode/bin/darknode && 
chmod +x ~/.darknode/bin/darknode && 
systemctl --user stop darknode &&
cp -a ~/.darknode/db/. ~/.darknode/db_bak/ &&
rm -rf ~/.darknode/db &&
mv ~/.darknode/db_bak ~/.darknode/db &&
echo %v > ~/.darknode/version &&
systemctl --user restart darknode`, url, ver)
	return util.RemoteRun(name, script)
}

func validateVersion(version string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := github.NewClient(nil)
	_, response, err := client.Repositories.GetReleaseByTag(ctx, "renproject", "darknode-release", version)
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
