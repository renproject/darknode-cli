package main

import (
	"io/ioutil"
	"encoding/json"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// Fixme : right now the bootstrap node address are hardcoded
var BootstrapNodes = func() []identity.MultiAddress {
	b1, _ := identity.NewMultiAddressFromString("/ip4/54.166.33.47/tcp/18514/republic/8MJ5XE5TwzknpizmVypzADLfYh5BfM")

	return []identity.MultiAddress{b1}
}()

// GetConfigOrGenerateNew will generate a new config for the darknode.
func GetConfigOrGenerateNew(directory , network string ) (config.Config, error) {

	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return config.Config{}, err
	}
	ethereumConfig := contract.Config{
		Network: contract.Network(network),
		URI:     "https://kovan.infura.io",
	}

	cfg:=  config.Config{
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

	configData, err := json.Marshal(cfg)
	if err != nil {
		return config.Config{}, err
	}
	if err := ioutil.WriteFile(directory + "/config.json", configData, 0600); err != nil {
		return config.Config{}, err
	}

	return cfg, nil
}
