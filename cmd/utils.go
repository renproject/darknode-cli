package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/crypto/ssh"
)

// Directory is the directory address of the cli and all darknodes data.
var Directory = path.Join(os.Getenv("HOME"), ".darknode")

// nodePath return the absolute directory of the node.
func nodePath(name string) string {
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

// getID parses the ID address from a bytes representation of
// multiAddress.
func getID(nodeDirectory string) (string, error) {
	addressFile := nodeDirectory + "/multiAddress.out"
	data, err := ioutil.ReadFile(addressFile)
	if err != nil {
		return "", err
	}
	multi, err := identity.NewMultiAddressFromString(strings.TrimSpace(string(data)))
	if err != nil {
		return "", err
	}

	return multi.ValueForProtocol(identity.RepublicCode)
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
		tagFile := path.Join(Directory, "darknodes", f.Name(), "tags.out")
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
	nodePath := nodePath(name)
	if _, err := os.Stat(nodePath); err != nil {
		return "", ErrNodeNotExist
	}
	if _, err := os.Stat(nodePath + "/config.json"); os.IsNotExist(err) {
		return "", ErrNodeNotExist
	}

	return nodePath, nil
}

// stringToEthereumAddress converts a hex string to a ethereum address.
// It returns an error if the provided string is an invalid address.
func stringToEthereumAddress(addr string) (common.Address, error) {
	if addr == "" {
		return common.Address{}, ErrEmptyAddress
	}
	if !common.IsHexAddress(addr) {
		return common.Address{}, ErrInvalidEthereumAddress
	}
	address := common.HexToAddress(addr)

	return address, nil
}

// run the command and pipe the output to the stdout
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// remoteRun runs the script to darknode with provided name.
func remoteRun(name, script string) error {
	// Parse the ssh private key
	key, err := ParsePrivateKey(name)
	if err != nil {
		return err
	}
	config := ssh.ClientConfig{
		User: "darknode",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to the instance using ssh
	ip, err := getIp(nodePath(name))
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", ip+":22", &config)
	if err != nil {
		return err
	}

	// Create a new session to run the script
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	return session.Run(script)
}

func redirectCommand() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "open", nil
	case "linux":
		return "xdg-open", nil
	default:
		return "", ErrUnsupportedOS
	}
}