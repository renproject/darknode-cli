package util

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

// GenerateSshKeyAndWriteToDir generates a new ssh key and write it to the given path.
func GenerateSshKeyAndWriteToDir(name string) error {
	path := NodePath(name)
	priKeyPath := filepath.Join(path, "ssh_keypair")
	pubKeyPath := filepath.Join(path, "ssh_keypair.pub")

	// Generate a random RSA key for ssh
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	key.Precompute()

	// Write the private key to file
	priKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priKeyBytes,
	}
	privatePEM := pem.EncodeToMemory(&privBlock)
	if err := ioutil.WriteFile(priKeyPath, privatePEM, 0600); err != nil {
		return err
	}

	// Write the public key to file
	publicRsaKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		return err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return ioutil.WriteFile(pubKeyPath, pubKeyBytes, 0600)
}

func ParseSshPrivateKey(name string) (ssh.Signer, error) {
	path := filepath.Join(NodePath(name), "ssh_keypair")
	sshKey, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(sshKey)
}

// EcdsaPrivateKeyToFile writes given private key to the target path in hex encoding.
func EcdsaPrivateKeyToFile(key *ecdsa.PrivateKey, path string) error {
	privateKeyBytes := crypto.FromECDSA(key)
	pkhex := hexutil.Encode(privateKeyBytes)
	pkhex = strings.Trim(pkhex, "0x")
	if len(pkhex) != 64 {
		return fmt.Errorf("invalid ecdsa key, expected 64 characters, got %v characters", len(pkhex))
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	return ioutil.WriteFile(path, []byte(pkhex), 0600)
}

// NewKeystoreFromECDSA parses a existing ecdsa key to a keystore.
func NewKeystoreFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *keystore.Key {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}
