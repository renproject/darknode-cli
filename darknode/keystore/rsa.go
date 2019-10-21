package keystore

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
)

// Rsa key for encrypting and decrypting sensitive data that must be transported
// between actors in the network.
type Rsa struct {
	*rsa.PrivateKey
}

// RandomRsaPrivKey using 2048 bits, with precomputed values for improved
// performance.
func RandomRsaPrivKey() (Rsa, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	privateKey.Precompute()
	if err != nil {
		return Rsa{}, err
	}
	return Rsa{
		PrivateKey: privateKey,
	}, nil
}

// MarshalJSON implements the json.Marshaler interface. The Rsa key is formatted
// according to the Ren Keystore specification.
func (key Rsa) MarshalJSON() ([]byte, error) {
	jsonKey := map[string]interface{}{}
	// Private key
	jsonKey["d"] = key.D.Bytes()
	jsonKey["primes"] = [][]byte{}
	for _, p := range key.Primes {
		jsonKey["primes"] = append(jsonKey["primes"].([][]byte), p.Bytes())
	}
	// Public key
	jsonKey["n"] = key.N.Bytes()
	jsonKey["e"] = key.E
	return json.Marshal(jsonKey)
}

// UnmarshalJSON implements the json.Unmarshaler interface. An Rsa key is
// created from data that is assumed to be compliant with the Ren Keystore
// specification. The Rsa key will be precomputed.
func (key *Rsa) UnmarshalJSON(data []byte) error {
	jsonKey := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &jsonKey); err != nil {
		return err
	}

	var err error

	// Private key
	key.PrivateKey = new(rsa.PrivateKey)
	key.PrivateKey.D, err = unmarshalBigIntFromMap(jsonKey, "d")
	if err != nil {
		return err
	}
	key.PrivateKey.Primes, err = unmarshalBigIntsFromMap(jsonKey, "primes")
	if err != nil {
		return err
	}

	// Public key
	key.PrivateKey.PublicKey = rsa.PublicKey{}
	key.PrivateKey.PublicKey.N, err = unmarshalBigIntFromMap(jsonKey, "n")
	if err != nil {
		return err
	}
	key.PrivateKey.PublicKey.E, err = unmarshalIntFromMap(jsonKey, "e")
	if err != nil {
		return err
	}

	key.Precompute()
	return nil
}
