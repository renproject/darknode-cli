package darknode

import (
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/renproject/aw"
	"github.com/renproject/darknode-cli/darknode/addr"
	"github.com/renproject/darknode-cli/darknode/keystore"
)

// Config is an in-memory description of the configuration file that will be
// loaded by a Darknode during boot.
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
	DNRAddress             common.Address `json:"dnrAddress"`
	ShifterRegistryAddress common.Address `json:"shifterRegistryAddress"`

	// Optional configuration
	HomeDir     *string         `json:"homeDir"`
	SentryDSN   *string         `json:"sentryDSN"`
	PeerOptions *aw.PeerOptions `json:"peerOptions"`
}

// NewConfig generate a new config.
func NewConfig(network Network) (Config, error) {
	// Generate a new random keystore.
	ks, err := keystore.RandomKeystore()
	if err != nil {
		return Config{}, err
	}

	// Parse the config or create a new random one
	return Config{
		Keystore:               ks,
		ECDSADistKeyShare:      ECDSADistKeyShare{PubKey: network.PublicKey()},
		Network:                network,
		Host:                   "0.0.0.0",
		Port:                   18514,
		Bootstraps:             network.BootstrapNodes(),
		DNRAddress:             network.DnrAddress(),
		ShifterRegistryAddress: network.ShiftRegistryAddress(),
		PeerOptions:            &aw.PeerOptions{
			DisablePeerDiscovery: true,
		},
	}, nil
}

// NewConfigFromJSONFile parses a json file that contains the config
// options specified by Config.
func NewConfigFromJSONFile(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	conf := Config{}
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return Config{}, err
	}

	if conf.HomeDir == nil {
		homeDir := filepath.Join(os.Getenv("HOME"), ".darknode")
		conf.HomeDir = &homeDir
	}
	if conf.PeerOptions == nil {
		conf.PeerOptions = &aw.PeerOptions{}
	}

	return conf, nil
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
