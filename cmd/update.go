package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/renproject/phi"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/urfave/cli"
)

// Location of the update script
const updateScriptLocation = "/home/darknode/.darknode/bin/update.sh"

// The command to execute to trigger an update. We use the -f to force an update.
var updateCommand = fmt.Sprintf("/bin/bash %s -f", updateScriptLocation)

// updateNode updates the Darknode to the latest release. It can also be used
// to update the config file of the darknode.
func updateNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	updateConfig := ctx.Bool("config")
	tags := ctx.String("tags")
	branch := ctx.String("branch")

	if tags == "" && name == "" {
		return ErrEmptyNodeName
	} else if tags == "" && name != "" {
		return updateSingleNode(name, branch, updateConfig)
	} else if tags != "" && name == "" {
		nodes, err := getNodesByTags(tags)
		if err != nil {
			return err
		}
		errs := make([]error, len(nodes))
		phi.ParForAll(nodes, func(i int) {
			errs[i] = updateSingleNode(nodes[i], branch, updateConfig)
		})
		return handleErrs(errs)
	}

	return ErrNameAndTags
}

func updateSingleNode(name, branch string, updateConfig bool) error {
	nodePath := nodePath(name)
	keyPairPath := path.Join(nodePath, "ssh_keypair")
	configPath := path.Join(nodePath, "config.json")
	ip, err := getIp(nodePath)
	if err != nil {
		return err
	}

	// Check if we need to update the node config
	if updateConfig {
		// Read the local config file
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return err
		}
		var cfg config.Config
		err = json.Unmarshal(data, &cfg)
		if err != nil {
			return err
		}

		// Read ssh private key from `ssh_keypair` file and decode it into a rsa key
		keyData, err := ioutil.ReadFile(keyPairPath)
		if err != nil {
			return err
		}
		block, _ := pem.Decode(keyData)
		if block == nil {
			return ErrInvalidSshKeyFile
		}
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}

		// Replace the rsaKey with the one for SSHing.
		cfg.Keystore.RsaKey = crypto.NewRsaKey(key)
		data, err = json.MarshalIndent(cfg, "", "    ")
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
			return err
		}

		updateConfigScript := fmt.Sprintf(`echo '%s' > $HOME/.darknode/config.json`, string(data))
		if err := run("ssh", "-i", keyPairPath, "darknode@"+ip, "-oStrictHostKeyChecking=no", updateConfigScript); err != nil {
			return err
		}

		fmt.Printf("%sConfig of [%s] has been updated to the local version.%s\n", GREEN, name, RESET)
	}

	err = run("ssh", "-i", keyPairPath, "darknode@"+ip, "-oStrictHostKeyChecking=no", updateCommand)
	if err != nil {
		return err
	}
	fmt.Printf("%s[%s] has been updated to the latest version.%s \n", GREEN, name, RESET)
	return nil
}

// sshNode will ssh into the Darknode
func sshNode(ctx *cli.Context) error {
	name := ctx.Args().First()
	if name == "" {
		cli.ShowCommandHelp(ctx, "ssh")
		return ErrEmptyNodeName
	}
	nodePath := nodePath(name)
	ip, err := getIp(nodePath)
	if err != nil {
		return err
	}
	keyPairPath := nodePath + "/ssh_keypair"

	return run("ssh", "-i", keyPairPath, "darknode@"+ip, "-oStrictHostKeyChecking=no")
}
