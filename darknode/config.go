package darknode

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/renproject/aw"
	"github.com/renproject/darknode-cli/darknode/addr"
	"github.com/renproject/darknode-cli/darknode/bindings"
	"github.com/renproject/darknode-cli/darknode/keystore"
)

// Config is an in-memory description of the configuration file that will be
// loaded by a Darknode during boot. Comparing to the `GeneralConfig` struct,
// this will always be the latest version of darknode config which we use to
// generate new config when deploying.
type Config struct {
	// Private configuration
	Keystore keystore.Keystore `json:"keystore"`

	// Public configuration
	Network    Network             `json:"network"`
	Host       string              `json:"host"`
	Port       int                 `json:"port"`
	Bootstraps addr.MultiAddresses `json:"bootstraps"`

	// Contract addresses
	ProtocolAddress common.Address `json:"protocolAddress"`

	// Optional configuration
	HomeDir       *string                `json:"homeDir"`
	SentryDSN     *string                `json:"sentryDSN"`
	PeerOptions   *aw.PeerOptions        `json:"peerOptions"`
	ClientOptions *aw.TCPConnPoolOptions `json:"clientOptions"`
	ServerOptions *aw.TCPServerOptions   `json:"serverOptions"`
}

// NewConfig generate a new config.
func NewConfig(network Network) (Config, error) {
	// Generate a new random keystore.
	ks, err := keystore.RandomKeystore()
	if err != nil {
		return Config{}, err
	}
	home := "/home/darknode/.darknode"

	// Parse the config or create a new random one
	return Config{
		Keystore:   ks,
		Network:    network,
		Host:       "0.0.0.0",
		Port:       18514,
		Bootstraps: network.BootstrapNodes(),

		ProtocolAddress: network.ProtocolAddr(),

		HomeDir:   &home,
		SentryDSN: nil,
		PeerOptions: &aw.PeerOptions{
			DisablePeerDiscovery: false,
		},
	}, nil
}

// NewConfigFromFile parses a json file that contains the config
// options specified by Config.
func NewConfigFromFile(filename string) (Config, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var conf Config
	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}

// ConfigToFile writes the Config to the target file in json format.
func ConfigToFile(config Config, path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "	")
	return encoder.Encode(config)
}

// GeneralConfig is the config struct which contains the common fields across
// all versions of darknode configs.
type GeneralConfig struct {
	// Private configuration
	Keystore keystore.Keystore `json:"keystore"`

	// Public configuration
	Network    Network             `json:"network"`
	Host       string              `json:"host"`
	Port       int                 `json:"port"`
	Bootstraps addr.MultiAddresses `json:"bootstraps"`

	// Contract addresses
	ProtocolAddress         common.Address `json:"protocolAddress"`
	DarknodeRegistryAddress common.Address `json:"dnrAddress"`
}

// DnrAddr returns the darknode registry contract address from the config. If
// only the protocol address is specified, it will try read the dnr address from
// the protocol contract.
// TODO : move to bindings or somewhere else
func (config GeneralConfig) DnrAddr(client *ethclient.Client) (common.Address, error) {
	if bytes.Equal(config.DarknodeRegistryAddress.Bytes(), common.Address{}.Bytes()) {
		protocol, err := bindings.NewProtocol(config.ProtocolAddress, client)
		if err != nil {
			return common.Address{}, err
		}
		return protocol.DarknodeRegistry(&bind.CallOpts{})
	}
	return config.DarknodeRegistryAddress, nil
}

// NewGeneralConfigFromFile parses a json file that contains the config
// options specified by GeneralConfig.
func NewGeneralConfigFromFile(filename string) (GeneralConfig, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return GeneralConfig{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return GeneralConfig{}, err
	}
	defer file.Close()

	var opts GeneralConfig
	err = json.NewDecoder(file).Decode(&opts)
	return opts, err
}
