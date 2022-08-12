package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/fatih/color"
	"github.com/renproject/darknode-cli/util"
	"github.com/urfave/cli"
)

// destroyNode tears down the deployed darknode by its name.
func destroyNode(ctx *cli.Context) error {
	force := ctx.Bool("force")
	name := ctx.Args().First()
	path := util.NodePath(name)
	if err := util.ValidateNodeExistence(name); err != nil {
		return err
	}

	if !force {
		// Last time confirm with user.
		text := "Are you sure you want to destroy your Darknode? (y/N)"
		input, err := util.Prompt(text)
		if err != nil {
			return err
		}
		input = strings.ToLower(strings.TrimSpace(text))
		if input != "yes" && input != "y" {
			return nil
		}
	}

	color.Green("Backing up config...")
	if err := util.BackUpConfig(name); err != nil {
		return err
	}

	color.Green("Destroying your Darknode...")
	destroy := fmt.Sprintf("cd %v && terraform destroy --force && cd .. && rm -rf %v", path, name)
	return util.Run("bash", "-c", destroy)
}

// export will export the darknode into a keystore file. The keystore file can
// be imported to wallets like metamask to withdraw unused funds.
func export(ctx *cli.Context) error {
	destination := ctx.String("out")
	name := ctx.Args().First()

	// Write to current directory if output directory not specified
	if destination == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		destination = filepath.Join(pwd, fmt.Sprintf("%v.json", name))
	}

	// Verify the node existence
	if err := util.ValidateNodeExistence(name); err != nil {
		return err
	}
	config, err := util.Config(name)
	if err != nil {
		panic(err)
	}

	// Ask fo for password
	text := "Please enter a password : \n"
	password, err := util.Prompt(text)
	if err != nil {
		return err
	}

	// Make sure user acknowledge the warnings
	text = color.YellowString("\n===========================================\n" +
		"This will export your darknode key into a keystore file.\n" +
		"You will need your keystore file + password to access your wallet.\n" +
		"Please save them in a secure location.\n" +
		"We CAN NOT retrieve or reset your keystore/password if you lose them.\n" +
		"Please acknowledge you have read and are aware of the above.(y/N)" +
		"\n===========================================\n")
	input, err := util.Prompt(text)
	if err != nil {
		return err
	}
	input = strings.ToLower(strings.TrimSpace(text))
	if input != "yes" && input != "y" {
		return nil
	}

	// Write to the keystore file
	key := util.NewKeystoreFromECDSA(config.Keystore.Ecdsa.PrivateKey)
	blob, err := keystore.EncryptKey(key, strings.TrimSpace(password), keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return err
	}
	if err := os.WriteFile(destination, blob, 0600); err != nil {
		return fmt.Errorf("failed to write keystore file: %v", err)
	}
	return nil
}
