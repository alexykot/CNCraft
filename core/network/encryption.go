package network

import (
	"crypto/cipher"
	"fmt"

	"github.com/alexykot/cncraft/pkg/protocol/crypto"
)

type crypter struct {
	enabled bool
	secret  []byte

	encrypt cipher.Stream
	decrypt cipher.Stream
}

func (c *crypter) Encrypt(data []byte) []byte {
	if !c.enabled {
		return data
	}

	output := make([]byte, len(data))
	c.encrypt.XORKeyStream(output, data)

	return output
}

func (c *crypter) Decrypt(data []byte) []byte {
	if !c.enabled {
		return data
	}

	output := make([]byte, len(data))
	c.decrypt.XORKeyStream(output, data)

	return output
}

func (c *crypter) Enable(secret []byte) error {
	encrypt, decrypt, err := crypto.NewEncryptAndDecrypt(secret)
	if err != nil {
		return fmt.Errorf("failed to enable conn encryption: %w", err)
	}

	c.enabled = true
	c.secret = secret
	c.encrypt = encrypt
	c.decrypt = decrypt
	return nil
}

func (c *crypter) Disable() {
	c.enabled = false
	c.secret = nil
	c.encrypt = nil
	c.decrypt = nil
}
