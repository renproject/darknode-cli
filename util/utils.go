package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/renproject/darknode-cli/darknode"
	"github.com/renproject/darknode-cli/darknode/addr"
	"golang.org/x/crypto/ssh"
)

var (
	// ErrEmptyNameAndTags is returned when both name and tags are not given.
	ErrEmptyNameAndTags = errors.New("please provide name or tags of the node you want to operate")

	// ErrTooManyArguments is returned when both name and tags are given.
	ErrTooManyArguments = errors.New("too many arguments, cannot have both name and tags")

	// ErrEmptyName is returned when user gives an empty node name.
	ErrEmptyName = errors.New("node name cannot be empty")

	// ErrUnknownDarknode is returned when the provided darknode name is unknown to us.
	ErrUnknownDarknode = errors.New("unknown darknode name")
)

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

func ValidateNodeName(name string) error {
	files, err := ioutil.ReadDir(filepath.Join(Directory, "/darknodes"))
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.Name() == name {
			return nil
		}
	}
	return ErrUnknownDarknode
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

// ID gets the ID of the node with given name.
func ID(name string) (addr.ID, error) {
	path := filepath.Join(NodePath(name), "config.json")
	config, err := darknode.NewConfigFromJSONFile(path)
	if err != nil {
		return addr.ID{}, err
	}
	return addr.FromPublicKey(config.Keystore.Ecdsa.PublicKey), nil
}

// IP gets the IP address of the node with given name.
func IP(name string) (string, error) {
	if name == "" {
		return "", ErrEmptyName
	}

	cmd := fmt.Sprintf("cd %v && terraform output ip", NodePath(name))
	ip, err := CommandOutput(cmd)
	return strings.TrimSpace(ip), err
}

func Network(name string) (darknode.Network, error) {
	path := filepath.Join(NodePath(name), "config.json")
	config, err := darknode.NewConfigFromJSONFile(path)
	if err != nil {
		return "", err
	}
	return config.Network, nil
}

// InitNodeDirectory creates the directory for the darknode.
func InitNodeDirectory(name, tags string) error {
	if name == "" {
		return ErrEmptyName
	}
	path := NodePath(name)

	// Ask user to destroy the old node first if there's already a node with the name.
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("Darknode [%v] already exists. \nIf you want to do a fresh deployment, destroy the old one first.", name)
	}

	// Make a directory for the new node
	if err := os.Mkdir(path, 0700); err != nil {
		return err
	}

	// Create the `tags.out` file if not exist.
	tagsPath := filepath.Join(path, "tags.out")
	if _, err := os.Stat(tagsPath); err != nil {
		return ioutil.WriteFile(tagsPath, []byte(strings.TrimSpace(tags)), 0600)
	}

	return nil
}

// run the command and pipe the output to the stdout
func Run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func SilentRun(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = ioutil.Discard
	cmd.Stderr = ioutil.Discard
	return cmd.Run()
}

func CommandOutput(commands string) (string, error) {
	cmd := exec.Command("bash", "-c", commands)
	output, err := cmd.Output()
	return string(output), err
}

// RemoteRun runs the script on the instance which host the darknode of given name.
func RemoteRun(name, script string) error {
	return RemoteRunWithUser(name, script, "darknode")
}

// RemoteRun runs the script on the instance as specific system user.
func RemoteRunWithUser(name, script, user string) error {
	// Parse the ssh private key
	key, err := ParseSshPrivateKey(name)
	if err != nil {
		return err
	}
	config := ssh.ClientConfig{
		User: user,
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

	// Redirect the remote stdin, stdout and stderr to local.
	sessStdIn, err := session.StdinPipe()
	if err != nil {
		return err
	}
	go io.Copy(sessStdIn, os.Stdin)
	sessStdOut, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stdout, sessStdOut)
	sessStdErr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go io.Copy(os.Stderr, sessStdErr)

	return session.Run(script)
}

func OpenInBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return Run("open", url)
	case "linux":
		return Run("xdg-open", url)
	}
	return nil
}

// RegisterUrl returns the url for registering a particular darknode.
func RegisterUrl(name string) (string, error) {
	path := filepath.Join(NodePath(name), "config.json")
	config, err := darknode.NewConfigFromJSONFile(path)
	if err != nil {
		return "", err
	}
	pubKey, err := ssh.NewPublicKey(&config.Keystore.Rsa.PublicKey)
	if err != nil {
		return "", err
	}
	pubKeyHex := hex.EncodeToString(pubKey.Marshal())
	id := addr.FromPublicKey(config.Keystore.Ecdsa.PublicKey)
	network, err := Network(name)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%v.renproject.io/darknode/%v?action=register&public_key=0x%s&name=%v", network, id.String(), pubKeyHex, name), nil
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
		if !haveAllTags {
			continue
		}

		// Check if the node is fully deployed
		if isDeployed(f.Name()) {
			nodes = append(nodes, f.Name())
		}
	}

	return nodes, nil
}

func isDeployed(name string) bool {
	path := NodePath(name)
	script := fmt.Sprintf("cd %v && terraform output ip", path)
	return SilentRun("bash", "-c", script) == nil
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

func Mkdir(path string, mode os.FileMode) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, mode)
	}
	return nil
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
