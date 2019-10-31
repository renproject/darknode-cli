package keystore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
)

// Ecdsa key for signing and verifying hashes.
type Ecdsa struct {
	*ecdsa.PrivateKey
}

// RandomEcdsaPrivKey using a secp256k1 s256 curve.
func RandomEcdsaPrivKey() (Ecdsa, error) {
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return Ecdsa{}, err
	}
	return Ecdsa{
		PrivateKey: privateKey,
	}, nil
}

// Address of the Ecdsa key. An Address is generated in the same way as an
// Ethereum address, but instead of a hex encoding, it uses a base58 encoding of
// a Keccak256 multihash.
func (key *Ecdsa) Address() string {
	bytes := elliptic.Marshal(secp256k1.S256(), key.PublicKey.X, key.PublicKey.Y)
	hash := crypto.Keccak256(bytes[1:]) // Keccak256 hash
	hash = hash[(len(hash) - 20):]      // Take the last 20 bytes
	addr := make([]byte, 2, 22)         // Create the multihash address
	addr[0] = multihash.KECCAK_256      // Set the keccak256 byte
	addr[1] = 20                        // Set the length byte
	addr = append(addr, hash...)        // Append the data
	return base58.EncodeAlphabet(addr, base58.BTCAlphabet)
}

// MarshalJSON implements the json.Marshaler interface. The Ecdsa key is
// formatted according to the Ren Keystore specification.
func (key Ecdsa) MarshalJSON() ([]byte, error) {
	jsonKey := map[string]interface{}{}
	// Private key
	jsonKey["d"] = key.D.Bytes()

	// Public key
	ethAddress, err := renAddressToEthAddress(key.Address())
	if err != nil {
		return []byte{}, err
	}

	jsonKey["address"] = ethAddress.Hex()
	jsonKey["x"] = key.X.Bytes()
	jsonKey["y"] = key.Y.Bytes()

	// Curve
	jsonKey["curveParams"] = map[string]interface{}{
		"p":    secp256k1.S256().P.Bytes(),  // the order of the underlying field
		"n":    secp256k1.S256().N.Bytes(),  // the order of the base point
		"b":    secp256k1.S256().B.Bytes(),  // the constant of the curve equation
		"x":    secp256k1.S256().Gx.Bytes(), // (x,y) of the base point
		"y":    secp256k1.S256().Gy.Bytes(),
		"bits": secp256k1.S256().BitSize, // the size of the underlying field
		"name": "s256",                   // the canonical name of the curve
	}
	return json.Marshal(jsonKey)
}

// UnmarshalJSON implements the json.Unmarshaler interface. An Ecdsa key is
// created from data that is assumed to be compliant with the Ren Keystore
// specification. The use of secp256k1 s256 curve is not checked.
func (key *Ecdsa) UnmarshalJSON(data []byte) error {
	jsonKey := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &jsonKey); err != nil {
		return err
	}

	var err error

	// Private key
	key.PrivateKey = new(ecdsa.PrivateKey)
	key.PrivateKey.D, err = unmarshalBigIntFromMap(jsonKey, "d")
	if err != nil {
		return err
	}

	// Public key
	key.PrivateKey.PublicKey = ecdsa.PublicKey{}
	key.PrivateKey.PublicKey.X, err = unmarshalBigIntFromMap(jsonKey, "x")
	if err != nil {
		return err
	}
	key.PrivateKey.PublicKey.Y, err = unmarshalBigIntFromMap(jsonKey, "y")
	if err != nil {
		return err
	}

	// Curve
	if jsonVal, ok := jsonKey["curveParams"]; ok {
		curveParams := elliptic.CurveParams{}
		jsonCurveParams := map[string]json.RawMessage{}
		if err := json.Unmarshal(jsonVal, &jsonCurveParams); err != nil {
			return err
		}
		curveParams.P, err = unmarshalBigIntFromMap(jsonCurveParams, "p")
		if err != nil {
			return err
		}
		curveParams.N, err = unmarshalBigIntFromMap(jsonCurveParams, "n")
		if err != nil {
			return err
		}
		curveParams.B, err = unmarshalBigIntFromMap(jsonCurveParams, "b")
		if err != nil {
			return err
		}
		curveParams.Gx, err = unmarshalBigIntFromMap(jsonCurveParams, "x")
		if err != nil {
			return err
		}
		curveParams.Gy, err = unmarshalBigIntFromMap(jsonCurveParams, "y")
		if err != nil {
			return err
		}
		curveParams.BitSize, err = unmarshalIntFromMap(jsonCurveParams, "bits")
		if err != nil {
			return err
		}
		curveParams.Name, err = unmarshalStringFromMap(jsonCurveParams, "name")
		if err != nil {
			return err
		}
		key.PrivateKey.Curve = &curveParams
	} else {
		return fmt.Errorf("curveParams is nil")
	}

	return nil
}

func renAddressToEthAddress(renAddr string) (common.Address, error) {
	bytes := base58.DecodeAlphabet(renAddr, base58.BTCAlphabet)[2:]
	if len(bytes) == 0 {
		return common.Address{}, errors.New("fail to decode the address")
	}
	ethAddr := common.BytesToAddress(bytes)
	return ethAddr, nil
}
