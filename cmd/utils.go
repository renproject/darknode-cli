package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/identity"
)

// Directory is the directory address of the cli and all darknodes data.
var Directory = path.Join(os.Getenv("HOME"), ".darknode")

// nodeDirectory return the absolute directory of the node.
func nodeDirectory(name string) string {
	return path.Join(Directory, "darknodes", name)
}

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

// getNodesByTags return the names of the nodes having the given tags.
func getNodesByTags(tags string) ([]string, error) {
	files, err := ioutil.ReadDir(Directory + "/darknodes")
	if err != nil {
		return nil, err
	}
	ts := strings.Split(strings.TrimSpace(tags), ",")
	nodes := make([]string, 0)
	for _, f := range files {
		tagFile := Directory + "/darknodes/" + f.Name() + "/tags.out"
		tags, err := ioutil.ReadFile(tagFile)
		if err != nil {
			continue
		}
		haveAllTags := true
		for i := range ts {
			if !strings.Contains(string(tags), ts[i]) {
				haveAllTags = false
			}
		}
		if haveAllTags {
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

// republicAddressToEthAddress converts republic address to ethereum address
func republicAddressToEthAddress(repAddress string) (common.Address, error) {
	addrByte := base58.DecodeAlphabet(repAddress, base58.BTCAlphabet)[2:]
	if len(addrByte) == 0 {
		return common.Address{}, errors.New("fail to decode the address")
	}

	address := common.BytesToAddress(addrByte)
	return address, nil
}

// copyFile copies the src file to dst. Any existing file will be overwritten
// and will not copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

// handleErrs checks a list of errors, return the first error encountered,
// nil otherwise.
func handleErrs(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// validateDarknodeName validates the darknode name and existence.
func validateDarknodeName(name string) (string, error) {
	if name == "" {
		return "", ErrEmptyNodeName
	}
	nodeDir := nodeDirectory(name)
	if _, err := os.Stat(nodeDir); err != nil {
		return "", ErrNodeNotExist
	}
	if _, err := os.Stat(nodeDir + "/config.json"); os.IsNotExist(err) {
		return "", ErrNodeNotExist
	}

	return nodeDir, nil
}