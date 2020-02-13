package util

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

	// ErrUnknownDarknode is returned when the provided darknode name is unknown to us.
	ErrUnknownDarknode = errors.New("unknown darknode name")
)

// ParseNodesFromNameAndTags returns the darknode names which satisfies the name
// requirements or the tag requirements.
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

// ValidateNodeName checks if there exists a node with given name.
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

// Config returns the config of the node with given name.
func Config(name string) (darknode.GeneralConfig, error){
	path := filepath.Join(NodePath(name), "config.json")
	return darknode.NewConfigFromJSONFile(path)
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
