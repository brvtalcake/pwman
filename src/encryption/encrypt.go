package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func Encrypt(pswd []byte, key []byte) ([]byte, error) {
	// convert password to bytes
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, pswd, nil)

	return ciphertext, nil
}
