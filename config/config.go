package config

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jbenet/go-base58"
)

type Keystore struct {
	Rsa   RsaKey   `json:"rsa"`
	Ecdsa EcdsaKey `json:"ecdsa"`
}

// EthereumConfig defines the network settings for Ethereum.
type EthereumConfig struct {
	Network string `json:"network"`
}

type Config struct {
	Keystore Keystore       `json:"keystore"`
	Address  string         `json:"address"`
	Ethereum EthereumConfig `json:"ethereum"`
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

	if conf.Ethereum.Network == "testnet" {
		conf.Ethereum.Network = "kovan"
	}
	return conf, nil
}

func (config *Config) EthereumAdress() (common.Address, error) {
	bytes := base58.DecodeAlphabet(config.Address, base58.BTCAlphabet)
	if len(bytes) < 22 {
		return [20]byte{}, errors.New("invalid address field in config")
	}
	bytes = bytes[2:]
	addrBytes := common.Address{}
	copy(addrBytes[:], bytes)
	return addrBytes, nil
}
