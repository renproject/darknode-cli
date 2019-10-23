package util

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/crypto/ssh"
)

// ErrEmptyNameAndTags is returned when both name and tags are not given.
var ErrEmptyNameAndTags = fmt.Errorf("%splease provide name or tags of the node you want to operate%s", RED, RESET)

// ErrTooManyArguments is returned when both name and tags are given.
var ErrTooManyArguments = fmt.Errorf("%stoo many arguments, cannot have both name and tags%s", RED, RESET)

// Directory is the directory address of the cli and all darknodes data.
var Directory = filepath.Join(os.Getenv("HOME"), ".darknode")

// NodePath return the absolute directory of the node with given name.
func NodePath(name string) string {
	return filepath.Join(Directory, "darknodes", name)
}

func ParseNodesFromNameAndTags(name, tags string) ([]string, error) {
	if name == "" && tags == "" {
		return nil, ErrEmptyNameAndTags
	} else if name == "" && tags != "" {
		return GetNodesByTags(tags)
	} else if name != "" && tags == "" {
		return []string{name}, nil
	} else {
		return nil, ErrTooManyArguments
	}
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

// IP parses the ip address from a bytes representation of multiAddress.
func IP(name string) (string, error) {
	// FIXME : instead read from the multiAddress file, ssh into the node and run command to get the ip address
	path := filepath.Join(NodePath(name), "multiAddress.out")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	multi, err := identity.NewMultiAddressFromString(strings.TrimSpace(string(data)))
	if err != nil {
		return "", err
	}

	return multi.ValueForProtocol(multiaddr.P_IP4)
}

// ID parses the ID address from a bytes representation of
// multiAddress.
func ID(name string) (string, error) {
	path := filepath.Join(NodePath(name), "multiAddress.out")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	multi, err := identity.NewMultiAddressFromString(strings.TrimSpace(string(data)))
	if err != nil {
		return "", err
	}

	return multi.ValueForProtocol(identity.RepublicCode)
}

// InitNodeDirectory creates the directory for the darknode.
func InitNodeDirectory(name, tags string) error {
	path := NodePath(name)

	// Make a directory for the new node if not exist.
	if _, err := os.Stat(path); err != nil {
		if err := os.Mkdir(path, 0700); err != nil {
			return err
		}
	}

	// Create the `tags.out` file if not exist.
	tagsPath := filepath.Join(path, "tags.out")
	if _, err := os.Stat(tagsPath); err != nil {
		return ioutil.WriteFile(tagsPath, []byte(strings.TrimSpace(tags)), 0644)
	}

	return nil
}

// run the command and pipe the output to the stdout
func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// RemoteRun runs the script on the instance which host the darknode of given name.
func RemoteRun(name, script string) error {
	return RemoteRunWithUser(name, script, "darknode")
}

// RemoteRun runs the script on the instance as specific system user.
func RemoteRunWithUser(name, script, user string) error {
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
	ip, err := IP(name)
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", ip), &config)
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

// GetNodesByTags return the names of the nodes which have the given tags.
func GetNodesByTags(tags string) ([]string, error) {
	files, err := ioutil.ReadDir(filepath.Join(Directory, "/darknodes"))
	if err != nil {
		return nil, err
	}
	ts := strings.Split(strings.TrimSpace(tags), ",")
	nodes := make([]string, 0)
	for _, f := range files {
		tagFile := filepath.Join(Directory, "darknodes", f.Name(), "tags.out")
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

// HandleErrs checks a list of errors, return the first error encountered,
// nil otherwise.
func HandleErrs(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// // validateDarknodeName validates the darknode name and existence.
// func validateDarknodeName(name string) (string, error) {
// 	if name == "" {
// 		return "", ErrEmptyNodeName
// 	}
// 	nodePath := nodePath(name)
// 	if _, err := os.Stat(nodePath); err != nil {
// 		return "", ErrNodeNotExist
// 	}
// 	if _, err := os.Stat(nodePath + "/config.json"); os.IsNotExist(err) {
// 		return "", ErrNodeNotExist
// 	}
//
// 	return nodePath, nil
// }
