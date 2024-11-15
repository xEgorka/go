// Package crypto provides interface and implements user data encryption.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Crypto describes crypto interface.
type Crypto interface {
	Key(usr, pass string) []byte
	Enc(plaintext string, key []byte) (string, error)
	Dec(ciphertext string, key []byte) (string, error)
}

// Crypto describes crypto engine.
type CryptoEngine struct{ timestamp time.Time }

// New creates new Crypto.
func New() Crypto { return &CryptoEngine{} }

const timeout = 5 // minutes since last access timestamp to crypto engine

// Key calculates user secret key.
func (c *CryptoEngine) Key(usr, pass string) []byte {
	k := sha256.Sum256([]byte(usr + pass))
	c.timestamp = time.Now()
	return k[:]
}

// Enc encrypts plaintext using key.
func (c *CryptoEngine) Enc(plaintext string, key []byte) (string, error) {
	if time.Now().Sub(c.timestamp).Minutes() > timeout {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	aes, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	c.timestamp = time.Now()
	return string(hex.EncodeToString(ciphertext)), nil
}

// Dec decrypts ciphertext using key.
func (c *CryptoEngine) Dec(ciphertext string, key []byte) (string, error) {
	if time.Now().Sub(c.timestamp).Minutes() > timeout {
		return "", status.Error(codes.Unauthenticated, "unauthenticated")
	}
	ct, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	aes, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ct := ct[:nonceSize], ct[nonceSize:]
	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ct), nil)
	if err != nil {
		return "", err
	}
	c.timestamp = time.Now()
	return string(plaintext), nil
}
