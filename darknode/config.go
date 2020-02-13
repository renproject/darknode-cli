package darknode

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"os"

	"github.com/btcsuite/btcd/btcec"
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
	Keystore          keystore.Keystore `json:"keystore"`
	ECDSADistKeyShare ECDSADistKeyShare `json:"ecdsaDistKeyShare"`

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
		Keystore:          ks,
		ECDSADistKeyShare: ECDSADistKeyShare{PubKey: network.PublicKey()},

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

// GeneralConfig is the config struct which contains the common fields across
// all versions of darknode configs.
type GeneralConfig struct {
	// Private configuration
	Keystore          keystore.Keystore `json:"keystore"`
	ECDSADistKeyShare ECDSADistKeyShare `json:"ecdsaDistKeyShare"`

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

// NewConfigFromJSONFile parses a json file that contains the config
// options specified by Config.
func NewConfigFromJSONFile(filename string) (GeneralConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return GeneralConfig{}, err
	}
	defer file.Close()

	var conf GeneralConfig
	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}

// The ECDSADistKeyShare is a temporary object used to store a Shamir's secret
// share of a ECDSA distributed key. Such a key is used by RenVM to sign
// transactions and messages as part of shifting tokens in/out of various
// distributed ledgers. In the future, it will be replaced by runtime storage so
// that there can be multiple ECDSA distributed keys that are constantly
// changed.
type ECDSADistKeyShare struct {
	PubKey       ecdsa.PublicKey `json:"pubKey"`
	PrivKeyShare []byte          `json:"privKeyShare,omitempty"`
}

// MarshalJSON implements the `json.Marshaler` interface for the
// ECDSADistKeyShare type.
func (dk ECDSADistKeyShare) MarshalJSON() ([]byte, error) {
	pubKey := map[string]interface{}{}
	pubKey["x"] = dk.PubKey.X
	pubKey["y"] = dk.PubKey.Y
	return json.Marshal(map[string]interface{}{
		"pubKey":       pubKey,
		"privKeyShare": dk.PrivKeyShare,
	})
}

// UnmarshalJSON implements the `json.Unmarshaler` interface for the
// ECDSADistKeyShare type.
func (dk *ECDSADistKeyShare) UnmarshalJSON(data []byte) (err error) {
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	pubKeyRaw := map[string]json.RawMessage{}
	if err = json.Unmarshal(m["pubKey"], &pubKeyRaw); err != nil {
		return
	}

	// Public key
	dk.PubKey.X = big.NewInt(0)
	if err = dk.PubKey.X.UnmarshalJSON(pubKeyRaw["x"]); err != nil {
		return
	}
	dk.PubKey.Y = big.NewInt(0)
	if err = dk.PubKey.Y.UnmarshalJSON(pubKeyRaw["y"]); err != nil {
		return
	}
	dk.PubKey.Curve = btcec.S256()

	// Private key share
	if err = json.Unmarshal(m["privKeyShare"], &dk.PrivKeyShare); err != nil {
		return
	}

	return nil
}
