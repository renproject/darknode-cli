package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/republicprotocol/republic-go/identity"
)

// StringInSlice checks whether the string is in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

// pipeToStd sets the input and output stream of the command to os standard
// input/output stream
func pipeToStd(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

// getIp parses the ip address from a bytes representation of
// multiAddress.
func getIp(nodeDirectory string) (string, error) {
	addressFile := nodeDirectory + "/multiAddress.out"
	data, err := ioutil.ReadFile(addressFile)
	if err != nil {
		return "", err
	}
	multi, err := identity.NewMultiAddressFromString(strings.TrimSpace(string(data)))
	if err != nil {
		return "", err
	}

	return multi.ValueForProtocol(identity.IP4Code)
}

// getNodesByTag return the names of the nodes having the given tag.
func getNodesByTag(tag string) ([]string, error) {
	files, err := ioutil.ReadDir(Directory + "/darknodes")
	if err != nil {
		return nil, err
	}
	nodes := []string{}

	for _, f := range files {
		tagFile := Directory + "/darknodes/" + f.Name() + "/tags.out"
		tags, err := ioutil.ReadFile(tagFile)
		if err != nil {
			continue
		}
		if strings.Contains(string(tags), tag) {
			nodes = append(nodes, f.Name())
		}
	}

	return nodes, nil
}

// cleanUp removes the directory
func cleanUp(nodeDirectory string) error {
	cleanCmd := exec.Command("rm", "-rf", nodeDirectory)
	if err := cleanCmd.Start(); err != nil {
		return err
	}

	return cleanCmd.Wait()
}
