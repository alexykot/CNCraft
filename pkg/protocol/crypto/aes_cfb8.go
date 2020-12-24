package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func NewEncryptAndDecrypt(secret []byte) (encrypt cipher.Stream, decrypt cipher.Stream, err error) {
	block, err := aes.NewCipher(secret)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to init aes cipherblock: %w", err)
	}
	if len(secret) != block.BlockSize() {
		return nil, nil,
			fmt.Errorf("aes cipherblock size %d does not match secret size %d", block.BlockSize(), len(secret))
	}

	encrypt = newCFB8(block, secret, false)
	decrypt = newCFB8(block, secret, true)

	return
}

func newCFB8(block cipher.Block, iv []byte, decrypt bool) cipher.Stream {
	blockSize := block.BlockSize()
	if len(iv) != blockSize {
		return nil
	}

	crypter := &crypterAESCFB8{
		b:       block,
		sr:      make([]byte, blockSize*4),
		srEnc:   make([]byte, blockSize),
		srPos:   0,
		decrypt: decrypt,
	}

	copy(crypter.sr, iv)

	return crypter
}

type crypterAESCFB8 struct {
	b       cipher.Block
	sr      []byte
	srEnc   []byte
	srPos   int
	decrypt bool
}

func (x *crypterAESCFB8) XORKeyStream(dst, src []byte) {
	blockSize := x.b.BlockSize()

	for i := 0; i < len(src); i++ {
		x.b.Encrypt(x.srEnc, x.sr[x.srPos:x.srPos+blockSize])

		var c byte
		if x.decrypt {
			c = src[i]
			dst[i] = c ^ x.srEnc[0]
		} else {
			c = src[i] ^ x.srEnc[0]
			dst[i] = c
		}

		x.sr[x.srPos+blockSize] = c
		x.srPos++

		if x.srPos+blockSize == len(x.sr) {
			copy(x.sr, x.sr[x.srPos:])
			x.srPos = 0
		}
	}
}
