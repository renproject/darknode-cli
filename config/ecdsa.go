package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
)

// ErrNilData is returned when a Verifier encounters a nil, or empty, data.
// Signing nil data is not in itself erroneous but it is rarely a reasonable
// action to take.
var ErrNilData = errors.New("nil data")

// ErrNilSignature is returned when a Verifier encounters a nil, or empty,
// signature.
var ErrNilSignature = errors.New("nil signature")

// ErrInvalidSignature is returned when a signature is invalid, or the
// recovered signatory does not match the required signatory.
var ErrInvalidSignature = errors.New("invalid signature")

// EcdsaKey for signing and verifying hashes.
type EcdsaKey struct {
	*ecdsa.PrivateKey
}

// RandomEcdsaKey using a secp256k1 s256 curve.
func RandomEcdsaKey() (EcdsaKey, error) {
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), rand.Reader)
	if err != nil {
		return EcdsaKey{}, err
	}
	return EcdsaKey{
		PrivateKey: privateKey,
	}, nil
}

// NewEcdsaKey returns an EcdsaKey from an existing private key. It does not
// verify that the private key was generated correctly.
func NewEcdsaKey(privateKey *ecdsa.PrivateKey) EcdsaKey {
	return EcdsaKey{
		PrivateKey: privateKey,
	}
}

// Sign implements the Signer interface. It uses the ecdsa.PrivateKey to sign
// the data without performing any kind of preprocessing of the data. If the
// data is not exactly 32 bytes, an error is returned.
func (key *EcdsaKey) Sign(data []byte) ([]byte, error) {
	return crypto.Sign(data, key.PrivateKey)
}

// Verify implements the Verifier interface. It uses its own address as the
// expected signatory.
func (key *EcdsaKey) Verify(data []byte, signature []byte) error {
	if data == nil || len(data) == 0 {
		return ErrNilData
	}
	if signature == nil || len(signature) == 0 {
		return ErrNilSignature
	}
	addrRecovered, err := RecoverAddress(data, signature)
	if err != nil {
		return err
	}
	if addrRecovered != key.Address() {
		return ErrInvalidSignature
	}
	return nil
}

// Address of the EcdsaKey. An Address is generated in the same way as an
// Ethereum address, but instead of a hex encoding, it uses a base58 encoding
// of a Keccak256 multihash.
func (key *EcdsaKey) Address() string {
	bytes := elliptic.Marshal(secp256k1.S256(), key.PublicKey.X, key.PublicKey.Y)
	hash := crypto.Keccak256(bytes[1:]) // Keccak256 hash
	hash = hash[(len(hash) - 20):]      // Take the last 20 bytes
	addr := make([]byte, 2, 22)         // Create the multihash address
	addr[0] = multihash.KECCAK_256      // Set the keccak256 byte
	addr[1] = 20                        // Set the length byte
	addr = append(addr, hash...)        // Append the data
	return base58.EncodeAlphabet(addr, base58.BTCAlphabet)
}

// Equal returns true if two EcdsaKeys are exactly equal. The name of the
// elliptic.Curve is not checked.
func (key *EcdsaKey) Equal(rhs *EcdsaKey) bool {
	return key.Address() == rhs.Address() &&
		key.D.Cmp(rhs.D) == 0 &&
		key.X.Cmp(rhs.X) == 0 &&
		key.Y.Cmp(rhs.Y) == 0 &&
		key.Curve.Params().P.Cmp(rhs.Curve.Params().P) == 0 &&
		key.Curve.Params().N.Cmp(rhs.Curve.Params().N) == 0 &&
		key.Curve.Params().B.Cmp(rhs.Curve.Params().B) == 0 &&
		key.Curve.Params().Gx.Cmp(rhs.Curve.Params().Gx) == 0 &&
		key.Curve.Params().Gy.Cmp(rhs.Curve.Params().Gy) == 0 &&
		key.Curve.Params().BitSize == rhs.Curve.Params().BitSize
}

// MarshalJSON implements the json.Marshaler interface. The EcdsaKey is
// formatted according to the Republic Protocol Keystore specification.
func (key EcdsaKey) MarshalJSON() ([]byte, error) {
	jsonKey := map[string]interface{}{}
	// Private key
	jsonKey["d"] = key.D.Bytes()

	// Public key
	ethAddress, err := republicAddressToEthAddress(key.Address())
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

// UnmarshalJSON implements the json.Unmarshaler interface. An EcdsaKey is
// created from data that is assumed to be compliant with the Republic Protocol
// Keystore specification. The use of secp256k1 s256 curve is not checked.
func (key *EcdsaKey) UnmarshalJSON(data []byte) error {
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

// Converts a republic address to an ethereum address
func republicAddressToEthAddress(repAddress string) (common.Address, error) {
	addByte := base58.DecodeAlphabet(repAddress, base58.BTCAlphabet)[2:]
	if len(addByte) == 0 {
		return common.Address{}, errors.New("fail to decode the address")
	}
	address := common.BytesToAddress(addByte)
	return address, nil
}

// RecoverAddress used to produce a signature.
func RecoverAddress(data []byte, signature []byte) (string, error) {
	// Returns 65-byte uncompress pubkey (0x04 | X | Y)
	publicKey, err := crypto.Ecrecover(data, signature)
	if err != nil {
		return "", err
	}

	// Address from an EcdsaKey
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
			X:     big.NewInt(0).SetBytes(publicKey[1:33]),
			Y:     big.NewInt(0).SetBytes(publicKey[33:65]),
		},
	}
	key := EcdsaKey{
		PrivateKey: privateKey,
	}
	return key.Address(), nil
}
