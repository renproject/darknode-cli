package main

import (
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// todo : any good way of defining the bootstrap node addresses
var BootstrapNodes = func() []identity.MultiAddress {
	b1, _ := identity.NewMultiAddressFromString("/ip4/34.203.9.146/tcp/18514/republic/8MJBssiB8aT6pGAM6MYj7YNUJTgxt7")
	b2, _ := identity.NewMultiAddressFromString("/ip4/54.250.246.106/tcp/18514/republic/8MJY6fvSCBCi3ujBqzTNTUkfF7WhFN")
	b3, _ := identity.NewMultiAddressFromString("/ip4/54.233.183.222/tcp/18514/republic/8MJ4LffVe6hDAha7AfKRt8Hr12xrVR")
	b4, _ := identity.NewMultiAddressFromString("/ip4/34.245.26.34/tcp/18514/republic/8MG9AZnq9s8UGqUcMMeq3r7azc58Mk")
	b5, _ := identity.NewMultiAddressFromString("/ip4/13.209.15.151/tcp/18514/republic/8MGLgu2wx8h5iiZDVLsgKLhTXqP7Uj")

	return []identity.MultiAddress{b1, b2, b3, b4, b5}
}()

// NewConfig will generate a new config for the darknode.
func NewConfig() (config.Config, error) {
	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return config.Config{}, err
	}
	ethereumConfig := ethereum.Config{
		Network: ethereum.NetworkKovan,
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

	return cfg, nil
}
