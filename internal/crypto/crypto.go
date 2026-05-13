package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

type Cipher struct {
	gcm cipher.AEAD
}

func NewCipher(key []byte) (*Cipher, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes long")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Cipher{gcm: gcm}, nil
}

func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to encrypt content: %w", err)
	}
	return c.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *Cipher) Decrypt(data []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	plaintext, err := c.gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	return plaintext, nil
}
