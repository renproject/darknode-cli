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

	if network != "testnet" && network != "falcon" && network != "nightly" {
		return config.Config{}, ErrUnknownProvider
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
							Path: "/home/ubuntu/.darknode/darknode.out",
						},
					},
				},
			},
			BootstrapMultiAddresses: BootstrapNodes(network),
			Ethereum: contract.Config{
				Network: contract.Network(network),
			},
		}
	} else {
		cfg, err = config.NewConfigFromJSONFile(configFile)
		if err != nil {
			return config.Config{}, nil
		}
	}

	configData, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return config.Config{}, err
	}
	if err := ioutil.WriteFile(directory+"/config.json", configData, 0600); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}

// Fixme : right now the bootstrap node address are hardcoded
func BootstrapNodes(network string) []identity.MultiAddress {
	addresess := make([]identity.MultiAddress, 0)

	switch network {
	case "testnet":
		b1, _ := identity.NewMultiAddressFromString("/ip4/54.250.246.106/tcp/18514/republic/8MJY6fvSCBCi3ujBqzTNTUkfF7WhFN")
		b2, _ := identity.NewMultiAddressFromString("/ip4/13.209.15.151/tcp/18514/republic/8MGLgu2wx8h5iiZDVLsgKLhTXqP7Uj")
		b3, _ := identity.NewMultiAddressFromString("/ip4/34.203.9.146/tcp/18514/republic/8MJBssiB8aT6pGAM6MYj7YNUJTgxt7")
		b4, _ := identity.NewMultiAddressFromString("/ip4/34.245.26.34/tcp/18514/republic/8MG9AZnq9s8UGqUcMMeq3r7azc58Mk")
		b5, _ := identity.NewMultiAddressFromString("/ip4/54.233.183.222/tcp/18514/republic/8MJ4LffVe6hDAha7AfKRt8Hr12xrVR")
		addresess = append(addresess, b1, b2, b3, b4, b5)
	case "falcon":
		b1, _ := identity.NewMultiAddressFromString("/ip4/13.124.184.167/tcp/18514/republic/8MJw8s6TVKmQH3kdM5kJUYqPmh3JmF")
		b2, _ := identity.NewMultiAddressFromString("/ip4/52.79.235.44/tcp/18514/republic/8MJEFcsQ5G8XMg5vka1XhswQotjbbj")
		b3, _ := identity.NewMultiAddressFromString("/ip4/13.114.234.59/tcp/18514/republic/8MKKUenZG8inoZqd8boYxGb1J3waAg")
		b4, _ := identity.NewMultiAddressFromString("/ip4/35.154.181.5/tcp/18514/republic/8MJNi8mgfQqD52bjUCnUzJ64uneJbk")
		addresess = append(addresess, b1, b2, b3, b4)
	case "nightly":
		b1, _ := identity.NewMultiAddressFromString("/ip4/54.255.182.246/tcp/18514/republic/8MHgRa2Uj7Tj2cgoA1PoULso7UmgVi")
		b2, _ := identity.NewMultiAddressFromString("/ip4/52.62.18.91/tcp/18514/republic/8MJnjSUVJCgP6YVjNWzaJXtiKE3p1o")
		addresess = append(addresess, b1, b2)
	default:
		panic(ErrUnknownNetwork)
	}

	return addresess
}
