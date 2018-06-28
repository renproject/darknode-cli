package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// Fixme : right now the bootstrap node address are hardcoded
var BootstrapNodes = func() []identity.MultiAddress {
	b1, _ := identity.NewMultiAddressFromString("/ip4/54.250.246.106/tcp/18514/republic/8MJY6fvSCBCi3ujBqzTNTUkfF7WhFN")
	b2, _ := identity.NewMultiAddressFromString("/ip4/13.209.15.151/tcp/18514/republic/8MGLgu2wx8h5iiZDVLsgKLhTXqP7Uj")
	b3, _ := identity.NewMultiAddressFromString("/ip4/34.203.9.146/tcp/18514/republic/8MJBssiB8aT6pGAM6MYj7YNUJTgxt7")
	b4, _ := identity.NewMultiAddressFromString("/ip4/34.245.26.34/tcp/18514/republic/8MG9AZnq9s8UGqUcMMeq3r7azc58Mk")
	b5, _ := identity.NewMultiAddressFromString("/ip4/54.233.183.222/tcp/18514/republic/8MJ4LffVe6hDAha7AfKRt8Hr12xrVR")

	return []identity.MultiAddress{b1, b2, b3, b4, b5}
}()

// GetConfigOrGenerateNew will generate a new config for the darknode.
func GetConfigOrGenerateNew(directory string) (config.Config, error) {

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return config.Config{}, err
	}
	ethereumConfig := contract.Config{
		Network: "testnet",
		URI:     "https://kovan.infura.io",
	}

	cfg := config.Config{
		Keystore:                keystore,
		Host:                    "0.0.0.0",
		Port:                    "18514",
		Address:                 identity.Address(keystore.Address()),
		BootstrapMultiAddresses: BootstrapNodes,
		Logs: logger.Options{
			Plugins: []logger.PluginOptions{
				{
					File: &logger.FilePluginOptions{
						Path: "/home/ubuntu/.darknode/darknode.out",
					},
				},
			},
		},
		Ethereum: ethereumConfig,
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
