package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/urfave/cli"
)

// GetConfigOrGenerateNew will generate a new config for the darknode.
func GetConfigOrGenerateNew(ctx *cli.Context, directory string) (config.Config, error) {
	keystoreFile := ctx.String("keystore")
	passphrase := ctx.String("passphrase")
	configFile := ctx.String("config")
	network := ctx.String("network")

	if network != "testnet" && network != "mainnet" {
		return config.Config{}, ErrUnknownNetwork
	}
	// Parse the keystore or create a new random one.
	var keystore crypto.Keystore
	var err error
	if keystoreFile == "" {
		keystore, err = crypto.RandomKeystore()
		if err != nil {
			return config.Config{}, err
		}
	} else {
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			return config.Config{}, err
		}
		if err := keystore.DecryptFromJSON(data, passphrase); err != nil {
			return config.Config{}, err
		}
	}

	// Parse the config or create a new random one
	var cfg config.Config
	if configFile == "" {
		cfg = config.Config{
			Keystore: keystore,
			Host:     "0.0.0.0",
			Port:     "18514",
			Address:  identity.Address(keystore.Address()),
			Logs: logger.Options{
				Plugins: []logger.PluginOptions{
					{
						File: &logger.FilePluginOptions{
							Path: "$HOME/.darknode/darknode.out",
						},
					},
				},
			},
			BootstrapMultiAddresses: BootstrapNodes(network),
			Ethereum: contract.Config{
				Network: contract.Network(network),
			},
			Alpha: 8,
		}
	} else {
		cfg, err = config.NewConfigFromJSONFile(configFile)
		if err != nil {
			return config.Config{}, nil
		}
	}

	// Write the config to file
	configData, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return config.Config{}, err
	}
	if err := ioutil.WriteFile(directory+"/config.json", configData, 0644); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

// Fixme : right now the bootstrap node address are hardcoded
func BootstrapNodes(network string) []identity.MultiAddress {
	addresess := make([]identity.MultiAddress, 0)

	switch network {
	case "mainnet":
		b1, _ := identity.NewMultiAddressFromString("/ip4/18.228.241.135/tcp/18514/republic/8MGdaR6EfrQdDb1MNDbYYvdx3KKqoc")
		// b2, _ := identity.NewMultiAddressFromString("/ip4/18.221.96.210/tcp/18514/republic/8MHdjyZxfuQT7unjYWzFCpm6m4qJNL")
		// b3, _ := identity.NewMultiAddressFromString("/ip4/52.9.204.195/tcp/18514/republic/8MG7prjbz51yqn1fqSygKSJfW3geBF")
		// b4, _ := identity.NewMultiAddressFromString("/ip4/54.201.216.97/tcp/18514/republic/8MKXcuQAjR2eEq8bsSHDPkYEmqmjtj")
		// b5, _ := identity.NewMultiAddressFromString("/ip4/52.77.88.84/tcp/18514/republic/8MJE7dUD8rHbYJ4RoWuVD6re4LKPVL")
		// b6, _ := identity.NewMultiAddressFromString("/ip4/54.252.152.19/tcp/18514/republic/8MJXyfjYuVZDPAfDm63G1NH1khgb2A")
		// b7, _ := identity.NewMultiAddressFromString("/ip4/18.195.208.147/tcp/18514/republic/8MHe93qDv1dBxoxygKLmeRwmMvpRam")
		// b8, _ := identity.NewMultiAddressFromString("/ip4/54.171.114.214/tcp/18514/republic/8MH2zgndxdCYReXL6uwggaCKGRHzdw")

		addresess = append(addresess, b1)
	case "testnet":
		b1, _ := identity.NewMultiAddressFromString("/ip4/18.211.224.194/tcp/18514/republic/8MJ7iKwcDxjndpD9EcXPgzKL9QJo2A")
		b2, _ := identity.NewMultiAddressFromString("/ip4/52.53.120.119/tcp/18514/republic/8MGdWRSn51Bc7ievAAkZ6x1hFAiJjf")
		b3, _ := identity.NewMultiAddressFromString("/ip4/52.53.120.119/tcp/18514/republic/8MGdWRSn51Bc7ievAAkZ6x1hFAiJjf")
		b4, _ := identity.NewMultiAddressFromString("/ip4/52.60.102.135/tcp/18514/republic/8MGATsfkjn1JtN9edB5C2kobUWtmJV")
		b5, _ := identity.NewMultiAddressFromString("/ip4/52.59.176.141/tcp/18514/republic/8MHX3awj1hj1x3XbKxq5uUuE8uQ6mq")
		b6, _ := identity.NewMultiAddressFromString("/ip4/18.228.50.197/tcp/18514/republic/8MJrKWdEmUFwZcKArzyABdYvovQQcP")

		addresess = append(addresess, b1, b2, b3, b4, b5, b6)
	default:
		panic(ErrUnknownNetwork)
	}

	return addresess
}
