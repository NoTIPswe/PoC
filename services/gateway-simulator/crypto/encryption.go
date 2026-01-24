package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type Encryptor struct {
	key   []byte
	keyID string
}

func NewEncryptor(hexKey string) (*Encryptor, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, errors.New("Invalid hex key")
	}
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	return &Encryptor{key: key, keyID: "static.key-v1"}, nil
}

func (e *Encryptor) KeyID() string {
	return e.keyID
}

func (e *Encryptor) Encrypt(plaintext []byte) (nonce, ciphertext []byte, err error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return nonce, ciphertext, nil
}
