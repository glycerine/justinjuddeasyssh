package easyssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

// CreateKeyPairFiles is the equivalent of running 'ssh-keygen -t rsa"'
func CreateKeyPairFiles(publicKeyPath, privateKeyPath string) error {

	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	return CreateKeyPair(publicKeyFile, privateKeyFile)
}

// CreateKeyPair creates a new SSH Key Pair writing the formatted keys to the corresponding io.Writers
func CreateKeyPair(publicKey, privateKey io.Writer) (err error) {
	k, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}
	privatePEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	err = pem.Encode(privateKey, privatePEM)
	if err != nil {
		return err
	}
	public, err := ssh.NewPublicKey(&k.PublicKey)
	if err != nil {
		return err
	}
	_, err = publicKey.Write(ssh.MarshalAuthorizedKey(public))
	return err
}

// LoadPrivateKey loads a file at the provided path and attempts to load it into an ssh.Signer that can be used for SSH servers
func LoadPrivateKey(filePath string) (ssh.Signer, error) {

	privateBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(privateBytes)
}
