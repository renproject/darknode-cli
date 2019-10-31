package addr

import (
	"crypto/ecdsa"
	cryptorand "crypto/rand"
	"errors"
	"math/rand"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/jbenet/go-base58"
	"github.com/multiformats/go-multihash"
)

// ErrInvalidID is returned when the address is invalid.
var ErrInvalidID = errors.New("invalid address")

// IDes is a wrapper around a slice of ID structs.
type IDes []ID

// Generate implements the `quick.Generator` interface. It returns realistic valid array of MultiAddress where the size of
// the array is determined by:`size`.
func (IDes) Generate(rand *rand.Rand, size int) reflect.Value {
	addrs := make(IDes, size)
	for i, addr := range addrs {
		addrs[i] = addr.Generate(rand, size).Interface().(ID)
	}
	return reflect.ValueOf(addrs)
}

// An ID identifies a Darknode in RenVM.
type ID struct {
	val string
}

// Generate implements the `quick.Generator` interface. It's used to create random values for testing, provides both
// valid and realistic values.
func (id ID) Generate(rand *rand.Rand, size int) reflect.Value {
	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), cryptorand.Reader)
	if err != nil {
		return reflect.ValueOf(err)
	}
	address := FromPublicKey(privateKey.PublicKey)
	return reflect.ValueOf(address)
}

// FromBase58 returns an ID parsed from a base58 string.
func FromBase58(val string) ID {
	return ID{val}
}

// FromEthereumAddress returns an ID parsed from an Ethereum address.
func FromEthereumAddress(ethaddress common.Address) ID {
	data := make([]byte, 2, 22)                // Create the multihash address
	data[0] = multihash.KECCAK_256             // Set the keccak256 byte
	data[1] = 20                               // Set the length byte
	data = append(data, ethaddress.Bytes()...) // Append the data
	return ID{base58.EncodeAlphabet(data, base58.BTCAlphabet)}
}

// FromPublicKey creates an ID from an ECDSA public key.
func FromPublicKey(pubKey ecdsa.PublicKey) ID {
	return FromEthereumAddress(crypto.PubkeyToAddress(pubKey))
}

// FromBytes returns a new ID parsed from the given bytes.
func FromBytes(raw [20]byte) ID {
	data := make([]byte, 2, 22)    // Create the multihash address
	data[0] = multihash.KECCAK_256 // Set the keccak256 byte
	data[1] = 20                   // Set the length byte
	data = append(data, raw[:]...) // Append the data
	return ID{base58.EncodeAlphabet(data, base58.BTCAlphabet)}
}

// ToBase58 returns the ID as a base58 encoded string.
func (id ID) ToBase58() string {
	return id.val
}

// ToEthereumAddress returns the Ethereum address representation of an ID.
func (id ID) ToEthereumAddress() (common.Address, error) {
	data := base58.DecodeAlphabet(id.ToBase58(), base58.BTCAlphabet)[2:]
	if len(data) < 20 {
		return common.Address{}, ErrInvalidID
	}
	return common.BytesToAddress(data), nil
}

// ToBytes returns the bytes representation of the REN address.
func (id ID) ToBytes() ([20]byte, error) {
	bytes := base58.DecodeAlphabet(id.ToBase58(), base58.BTCAlphabet)
	if len(bytes) < 22 {
		return [20]byte{}, ErrInvalidID
	}
	addrBytes := [20]byte{}
	copy(addrBytes[:], bytes[2:])
	return addrBytes, nil
}

// Equal compares two IDes. If the two IDes are equal, then it returns
// true. Otherwise, it returns false.
func (id ID) Equal(other ID) bool {
	return id.val == other.val
}

// String implements the `fmt.Stringer` interface.
func (id ID) String() string {
	return id.ToBase58()
}
