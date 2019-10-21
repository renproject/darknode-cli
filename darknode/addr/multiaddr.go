package addr

import (
	"bytes"
	"crypto/ecdsa"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multihash"
	"golang.org/x/crypto/sha3"
)

// ErrMalformedMultiAddressLength indicates that the length of the multi-address data buffer is invalid.
var ErrMalformedMultiAddressLength = errors.New("malformed multi-address length")

// Codes for the multi-address protocol. See https://multiformats.io/multiaddr for more information.
const (
	IP4Code = 0x0004
	IP6Code = 0x0029
	TCPCode = 0x0006
	RenCode = 0x0065
)

func init() {
	// Register then Ren code during package initalization.
	multiaddr.AddProtocol(multiaddr.Protocol{
		Code:       RenCode,
		Size:       multiaddr.LengthPrefixedVarSize,
		Name:       "ren",
		Path:       false,
		Transcoder: multiaddr.NewTranscoderFromFunctions(stB, btS, nil),
	})
}

type MultiAddresses []MultiAddress

// Generate implements the `quick.Generator` interface. It returns realistic valid array of MultiAddress where the size of
// the array is determined by:`size`.
func (MultiAddresses) Generate(rand *rand.Rand, size int) reflect.Value {
	addrs := make(MultiAddresses, size)
	for i, addr := range addrs {
		addrs[i] = addr.Generate(rand, size).Interface().(MultiAddress)
	}
	return reflect.ValueOf(addrs)
}

type MultiAddress struct {
	value     multiaddr.Multiaddr
	nonce     uint64
	signature [65]byte
}

// Generate implements the `quick.Generator` interface. It's used to create random values for testing, provides both
// valid and realistic values.
func (multi MultiAddress) Generate(rand *rand.Rand, size int) reflect.Value {
	ip4 := fmt.Sprintf("%v.%v.%v.%v", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
	tcp := fmt.Sprintf("%v", rand.Intn(8000))

	privateKey, err := ecdsa.GenerateKey(secp256k1.S256(), cryptorand.Reader)
	if err != nil {
		return reflect.ValueOf(err)
	}
	address := FromPublicKey(privateKey.PublicKey)

	value := fmt.Sprintf("/ip4/%v/tcp/%v/ren/%s", ip4, tcp, address.ToBase58())
	multiAddr, err := NewMultiAddressFromString(value)
	if err != nil {
		return reflect.ValueOf(err)
	}
	if err := multiAddr.Sign(privateKey); err != nil {
		return reflect.ValueOf(err)
	}

	return reflect.ValueOf(multiAddr)
}

// NewMultiAddressFromString returns a new MultiAddress that uses the current
// seconds since UNIX epoch as its nonce, but is unsigned. An error is returned
// if the string is not in the form
// `"/ip4/{ip-address}/tcp/{port}/ren/{base58-encoded-id}"`.
func NewMultiAddressFromString(value string) (MultiAddress, error) {
	multiAddr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return MultiAddress{}, err
	}

	_, err = multiAddr.ValueForProtocol(RenCode)
	if err != nil {
		return MultiAddress{}, err
	}
	ip4Addr, err := multiAddr.ValueForProtocol(IP4Code)
	if err != nil {
		return MultiAddress{}, err
	}
	port, err := multiAddr.ValueForProtocol(TCPCode)
	if err != nil {
		return MultiAddress{}, err
	}
	if _, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", ip4Addr, port)); err != nil {
		return MultiAddress{}, err
	}

	return MultiAddress{
		value: multiAddr,
		nonce: uint64(time.Now().Unix()),
	}, err
}

func NewSignedMultiAddressFromString(value, signature string) (MultiAddress, error) {
	multiAddr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return MultiAddress{}, err
	}

	_, err = multiAddr.ValueForProtocol(RenCode)
	if err != nil {
		return MultiAddress{}, err
	}
	ip4Addr, err := multiAddr.ValueForProtocol(IP4Code)
	if err != nil {
		return MultiAddress{}, err
	}
	port, err := multiAddr.ValueForProtocol(TCPCode)
	if err != nil {
		return MultiAddress{}, err
	}
	if _, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", ip4Addr, port)); err != nil {
		return MultiAddress{}, err
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return MultiAddress{}, err
	}
	var sig [65]byte
	copy(sig[:], sigBytes)

	return MultiAddress{
		value:     multiAddr,
		nonce:     uint64(time.Now().Unix()),
		signature: sig,
	}, err
}

func multiAddressFromArgs(value string, nonce uint64, signature [65]byte) (MultiAddress, error) {
	multiAddress, err := NewMultiAddressFromString(value)
	if err != nil {
		return MultiAddress{}, err
	}
	multiAddress.nonce = nonce
	multiAddress.signature = signature
	return multiAddress, nil
}

func (multiAddr MultiAddress) String() string {
	return multiAddr.value.String()
}

func (multiAddr MultiAddress) Equal(other MultiAddress) bool {
	return multiAddr.value.Equal(other.value) &&
		multiAddr.nonce == other.nonce &&
		bytes.Equal(multiAddr.signature[:], other.signature[:])
}

func (multiAddr MultiAddress) ID() ID {
	ren, err := multiAddr.value.ValueForProtocol(RenCode)
	if err != nil {
		panic(err)
	}
	return FromBase58(ren)
}

func (multiAddr MultiAddress) Nonce() uint64 {
	return multiAddr.nonce
}

func (multiAddr MultiAddress) IP4() string {
	ip4, err := multiAddr.value.ValueForProtocol(IP4Code)
	if err != nil {
		panic(err)
	}
	return ip4
}

func (multiAddr MultiAddress) Port() int {
	portString, err := multiAddr.value.ValueForProtocol(TCPCode)
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic(err)
	}
	return port
}

func (multiAddr MultiAddress) NetworkAddress() net.Addr {
	ipAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", multiAddr.IP4(), multiAddr.Port()))
	if err != nil {
		panic(err)
	}
	return ipAddr
}

func (multiAddr MultiAddress) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, multiAddr.signature); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, multiAddr.Nonce()); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, []byte(multiAddr.String())); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (multiAddr *MultiAddress) UnmarshalBinary(data []byte) error {
	if len(data) < 74 {
		return ErrMalformedMultiAddressLength
	}
	buf := bytes.NewReader(data[:73])

	sig := [65]byte{}
	if err := binary.Read(buf, binary.BigEndian, &sig); err != nil {
		return err
	}

	nonce := uint64(0)
	if err := binary.Read(buf, binary.BigEndian, &nonce); err != nil {
		return err
	}

	umarshaledMultiAddr, err := multiAddressFromArgs(string(data[73:]), nonce, sig)
	if err != nil {
		return err
	}

	*multiAddr = umarshaledMultiAddr
	return nil
}

type multiAddressJsonValue struct {
	Value     string `json:"value"`
	Nonce     uint64 `json:"nonce"`
	Signature []byte `json:"signature"`
}

// MarshalJSON implements the json.Marshaler interface.
func (multiAddr *MultiAddress) MarshalJSON() ([]byte, error) {
	if multiAddr.value == nil {
		return []byte{}, errors.New("value cannot be nil")
	}
	val := multiAddressJsonValue{
		Signature: multiAddr.signature[:],
		Nonce:     multiAddr.nonce,
		Value:     multiAddr.value.String(),
	}
	return json.Marshal(val)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (multiAddr *MultiAddress) UnmarshalJSON(data []byte) error {
	val := multiAddressJsonValue{}
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	sig := [65]byte{}
	copy(sig[:], val.Signature)
	newMultiAddr, err := multiAddressFromArgs(val.Value, val.Nonce, sig)
	if err != nil {
		return err
	}
	*multiAddr = newMultiAddr
	return nil
}

func (multiaddr *MultiAddress) Sign(key *ecdsa.PrivateKey) error {
	data := []byte(multiaddr.String())

	hashSum256 := sha3.Sum256(data)
	signature, err := crypto.Sign(hashSum256[:], key)
	if err != nil {
		return err
	}
	copy(multiaddr.signature[:], signature)

	return nil
}

func (multiaddr *MultiAddress) Verify() bool {
	data := []byte(multiaddr.String())
	hashSum256 := sha3.Sum256(data)
	pubKey, err := crypto.SigToPub(hashSum256[:], multiaddr.signature[:])
	if err != nil {
		return false
	}

	return FromPublicKey(*pubKey) == multiaddr.ID()
}

func (multiaddr *MultiAddress) PublicKey() (*ecdsa.PublicKey, error) {
	data := []byte(multiaddr.String())
	hashSum256 := sha3.Sum256(data)
	return crypto.SigToPub(hashSum256[:], multiaddr.signature[:])
}

// Convert string to bytes.
func stB(s string) ([]byte, error) {
	m, err := multihash.FromB58String(s)
	if err != nil {
		return nil, fmt.Errorf("error parsing Ren address %s: %s", s, err)
	}
	size := multiaddr.CodeToVarint(len(m))
	b := append(size, m...)
	return b, nil
}

// Convert bytes to string.
func btS(b []byte) (string, error) {
	size, n, err := multiaddr.ReadVarintCode(b)
	if err != nil {
		return "", err
	}
	b = b[n:]
	if len(b) != size {
		return "", errors.New("malformed length")
	}
	m, err := multihash.Cast(b)
	if err != nil {
		return "", err
	}
	return m.B58String(), nil
}
