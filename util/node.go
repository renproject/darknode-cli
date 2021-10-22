package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
)

// ParseNodesFromNameAndTags returns the darknode names which satisfies the name
// requirements or the tag requirements.
func ParseNodesFromNameAndTags(name, tags string) ([]string, error) {
	if name == "" && tags == "" {
		return nil, ErrEmptyNameAndTags
	} else if name == "" && tags != "" {
		return GetNodesByTags(tags)
	} else if name != "" && tags == "" {
		return []string{name}, NodeExistence(name)
	} else {
		return nil, ErrTooManyArguments
	}
}

// ValidateName validates the given darknode name. It should
// 1) Only contains letter, number, "-" and "_".
// 2) No more than 32 characters
// 3) Do not contain any whitespace
func ValidateName(name string) error {
	if strings.TrimSpace(name) != name {
		return fmt.Errorf("name cannot have whitespace")
	}

	nameRegex, err := regexp.Compile("^[a-zA-Z0-9_-]{1,32}$")
	if err != nil {
		return err
	}
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("darknode name should be less than 32 characters and not contain any special character")
	}
	return nil
}

// NodeExistence checks if there exists a node with given name.
func NodeExistence(name string) error {
	path := filepath.Join(Directory, "darknodes", name)
	_, err := os.Stat(path)
	return err
}

// Config returns the config of the node with given name.
func Config(name string) (darknode.GeneralConfig, error) {
	path := filepath.Join(NodePath(name), "config.json")
	return darknode.NewGeneralConfigFromJSONFile(path)
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
	if strings.HasPrefix(ip, "\"") {
		ip = strings.Trim(ip, "\"")
	}
	return strings.TrimSpace(ip), err
}

// Version gets the version of the software the darknode currently is running.
func Version(name string) string {
	script := "cat ~/.darknode/version"
	version, err := RemoteOutput(name, script)
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(version))
}

// Network gets the network of the darknode.
func Network(name string) (darknode.Network, error) {
	path := filepath.Join(NodePath(name), "config.json")
	config, err := darknode.NewConfigFromJSONFile(path)
	if err != nil {
		return "", err
	}
	return config.Network, nil
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
	nodes := make([]string, 0)
	for _, f := range files {
		path := filepath.Join(Directory, "darknodes", f.Name(), "tags.out")
		tagFile, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}
		if !ValidateTags(string(tagFile), tags) {
			continue
		}

		// Check if the node is fully deployed
		if isDeployed(f.Name()) {
			nodes = append(nodes, f.Name())
		}
	}
	if len(nodes) == 0 {
		return nil, errors.New("cannot find any darknode with given tags")
	}

	return nodes, nil
}

func ValidateTags(have, required string) bool {
	tagsStr := strings.Split(strings.TrimSpace(required), ",")
	for _, tag := range tagsStr {
		if !strings.Contains(have, tag) {
			return false
		}
	}
	return true
}

func isDeployed(name string) bool {
	path := NodePath(name)
	script := fmt.Sprintf("cd %v && terraform output ip", path)
	return SilentRun("bash", "-c", script) == nil
}
