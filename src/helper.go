package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/republicprotocol/republic-go/identity"
)

// ErrNoNodesFound is returned when no nodes can be found with the given tag.
var ErrNoNodesFound = fmt.Errorf("%sno nodes can be found with the given tag%s", red, reset)

// StringInSlice checks whether the string is in the slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Mkdir tries to create folder with the given name. If not success,
// it will try appending number appended to the name.
func Mkdir(name string) (string, error) {
	// Try creating directory with the exact name we get
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err = os.Mkdir(name, os.ModePerm)
		if err == nil {
			return name, nil
		}
	}

	// Try creating directory with number appended to the name.
	i := 1
	for {
		dirName := fmt.Sprintf("%v_%d", name, i)
		if _, err := os.Stat(dirName); os.IsNotExist(err) {
			err = os.Mkdir(dirName, os.ModePerm)
			if err != nil {
				return "", err
			}
			return dirName, nil
		}
		i++
	}
}

// pipeToStd set the input and output stream of the command to os standard
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
		if strings.Contains(string(tags), strings.ToLower(tag)) {
			nodes = append(nodes, f.Name())
		}
	}

	return nodes, nil
}
