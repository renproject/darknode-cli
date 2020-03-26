package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/renproject/darknode-cli/util"
	"github.com/renproject/phi"
	"github.com/urfave/cli"
)

// updateNode updates the Darknode to the latest release. It can also be used
// to update the config file of the darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	tags := ctx.String("tags")
	version := ctx.String("version")

	nodes, err := util.ParseNodesFromNameAndTags(name, tags)
	if err != nil {
		return err
	}
	if !compareVersion(nodes, version){
		return nil
	}

	errs := make([]error, len(nodes))
	phi.ParForAll(nodes, func(i int) {
		errs[i] = updateSingleNode(nodes[i], version)
	})
	return util.HandleErrs(errs)
}

// Compare the current version and target version. Show warning to user if the
// current version is greater than the target version and confirm with user
// before updating. Returned boolean indicates whether we want to continue.
func compareVersion(nodes []string, new string) bool {
	curVersions := make([]string, len(nodes))
	phi.ParForAll(nodes, func(i int) {
		// Ignore the error. Leave the version as empty string if we cannot get its current version.
		curVersions[i], _ = getCurrentVersion(nodes[i])
	})

	newVersion, err := version.NewVersion(new)
	if err != nil {
		return true
	}
	lower := false
	for i, ver := range curVersions{
		if ver == ""{
			continue
		}
		curVersion , err := version.NewVersion(ver)
		if err != nil {
			continue
		}
		if curVersion.GreaterThan(newVersion){
			color.Red("Darknode [%v] currently running with version %v, and you are trying to downgrade to %v", nodes[i], curVersion.String(), newVersion.String())
			lower = true
		}
	}
	if lower {
		color.Red("Do you want to continue?(y/N)")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		input := strings.ToLower(strings.TrimSpace(text))
		return input == "yes" || input == "y"
	}
	return true
}

func updateSingleNode(name, ver string) error {
	url := "https://www.github.com/renproject/darknode-release/releases/latest/download/darknode"
	if ver != "" {
		if err := validateVersion(ver); err != nil {
			return err
		}
		url = fmt.Sprintf("https://github.com/renproject/darknode-release/releases/download/%v/darknode", ver)
	}

	script := fmt.Sprintf(`mv ~/.darknode/bin/darknode ~/.darknode/bin/darknode-backup && curl -sL %v > ~/.darknode/bin/darknode && chmod +x ~/.darknode/bin/darknode && systemctl --user restart darknode`, url)
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

func getCurrentVersion(name string) (string, error) {
	request := `{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "ren_queryStat",
    "params": {}
}`
	buf := bytes.NewBuffer([]byte(request))
	ip, err:= util.IP(name)
	if err != nil {
		return"",  err
	}
	url := fmt.Sprintf("http://%v:18515", ip)
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()
	req, err:= http.NewRequestWithContext(ctx, "POST", url, buf)
	if err != nil {
		return"",  err
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err!= nil {
		return "", err
	}
	response := struct{
		Result struct{
			Version string `json:"version"`
		} `json:"result"`
	} {}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return"",  err
	}
	versionString := strings.Split(response.Result.Version, "-")
	if len(versionString) < 1 {
		return "", fmt.Errorf("invalid version = %v", response.Result.Version)
	}
	return versionString[0], nil
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
