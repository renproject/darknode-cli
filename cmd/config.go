package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/renproject/aw"
	"github.com/republicprotocol/darknode-cli/darknode"
	"github.com/republicprotocol/darknode-cli/darknode/keystore"
	"github.com/urfave/cli"
)

// GetConfigOrGenerateNew will generate a new config for the darknode.
func GetConfigOrGenerateNew(ctx *cli.Context, directory string) (darknode.Config, error) {
	// Parse the network of the darknode
	network, err := darknode.NewNetwork(ctx.String("network"))
	if err != nil {
		return darknode.Config{}, err
	}

	// Generate a new random keystore.
	ks, err := keystore.RandomKeystore()
	if err != nil {
		return darknode.Config{}, err
	}

	// Parse the config or create a new random one
	config := darknode.Config{
		Keystore:               ks,
		ECDSADistKeyShare:      darknode.ECDSADistKeyShare{PubKey: network.PublicKey()},
		Network:                network,
		Host:                   "0.0.0.0",
		Port:                   18514,
		Bootstraps:             network.BootstrapNodes(),
		DNRAddress:             network.DnrAddress(),
		ShifterRegistryAddress: network.ShiftRegistryAddress(),
		PeerOptions:            &aw.PeerOptions{},
	}

	// Write the config to file
	configData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return darknode.Config{}, err
	}
	configPath := filepath.Join(directory, "/config.json")
	err = ioutil.WriteFile(configPath, configData, 0600)
	return config, err
}
