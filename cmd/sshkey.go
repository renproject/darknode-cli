package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// WriteSshKey writes the rsa key to a file in the ssh format.
func WriteSshKey(rsaKey *rsa.PrivateKey, directory string) error {
	// Path to save the ssh keys
	keyPairPath := directory + "/ssh_keypair"
	pubKeyPath := directory + "/ssh_keypair.pub"

	// Write the private key to file
	priKeyBytes := x509.MarshalPKCS1PrivateKey(rsaKey)
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
		return err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	return ioutil.WriteFile(pubKeyPath, pubKeyBytes, 0600)
}

// StringfySshPubkey returned the
func StringfySshPubkey(key ssh.PublicKey) string {
	pubKeyBytes := ssh.MarshalAuthorizedKey(key)

	return string(pubKeyBytes)
}
