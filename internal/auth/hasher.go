package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type hasher struct {
	memory      uint32
	saltLength  uint32
	time        uint32
	keyLen      uint32
	parallelism uint8
}

func NewHasher() *hasher {
	return &hasher{
		memory:      64 * 1024,
		saltLength:  16,
		parallelism: 4,
		time:        1,
		keyLen:      32,
	}
}

func (h *hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.time,
		h.memory,
		h.parallelism,
		h.keyLen,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.memory, h.time, h.parallelism, b64Salt, b64Hash,
	)

	return encoded, nil
}
