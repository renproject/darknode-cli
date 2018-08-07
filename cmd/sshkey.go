package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/republicprotocol/republic-go/crypto"
	"golang.org/x/crypto/ssh"
)

// NewSshKeyPair generate a new ssh key pair and writes the keys into files.
// It returns the public ssh key and the path of the rsa key file.
func NewSshKeyPair(directory string) (ssh.PublicKey, error) {
	// Path to save the ssh keys
	keyPairPath := directory + "/ssh_keypair"
	pubKeyPath := directory + "/ssh_keypair.pub"

	rsaKey, err := crypto.RandomRsaKey()
	if err != nil {
		return nil, nil
	}

	// Write the private key to file
	priKeyBytes := x509.MarshalPKCS1PrivateKey(rsaKey.PrivateKey)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priKeyBytes,
	}
	privatePEM := pem.EncodeToMemory(&privBlock)
	ioutil.WriteFile(keyPairPath, privatePEM, 0600)

	// Write the public key to file
	publicRsaKey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	err = ioutil.WriteFile(pubKeyPath, pubKeyBytes, 0600)

	return publicRsaKey, err

}

// StringfySshPubkey returned the
func StringfySshPubkey(key ssh.PublicKey) string {
	pubKeyBytes := ssh.MarshalAuthorizedKey(key)

	return string(pubKeyBytes)
}
