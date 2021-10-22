package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"path/filepath"

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
