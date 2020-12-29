package mojang

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
)

type RSACrypter struct {
	privateKey      *rsa.PrivateKey
	publicKey       *rsa.PublicKey
	privateKeyBytes []byte
	publicKeyBytes  []byte
}

func NewRSACrypter() (*RSACrypter, error) {
	privateKey, publicKey, err := generateRandomKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create keys: %w", err)
	}

	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyDER, _ := x509.MarshalPKIXPublicKey(publicKey)

	return &RSACrypter{
		privateKey:      privateKey,
		publicKey:       publicKey,
		privateKeyBytes: privateKeyDER,
		publicKeyBytes:  publicKeyDER,
	}, nil
}

func (c *RSACrypter) GetPubKey() []byte {
	return c.publicKeyBytes
}

func (c *RSACrypter) Encrypt(data []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, data)
}

func (c *RSACrypter) Decrypt(data []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, c.privateKey, data)
}

func generateRandomKey() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate new RSA key: %w", err)
	}

	privateKey.Precompute()
	if err := privateKey.Validate(); err != nil {
		return nil, nil, fmt.Errorf("failed to validate new RSA key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}
