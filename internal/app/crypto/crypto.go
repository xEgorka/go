// Package crypto provides interface and implements encryption.
package crypto

import "golang.org/x/crypto/bcrypt"

// Crypto describes crypto interface.
type Crypto interface {
	Hash(pass string) (string, error)
	Verify(pass, hash string) error
}

// CryptoEngine describes crypto engine.
type CryptoEngine struct{}

// New creates new Crypto.
func New() Crypto { return &CryptoEngine{} }

const salt = "salt"

// Hash calculates user password hash.
func (c *CryptoEngine) Hash(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass+salt), bcrypt.DefaultCost)
	if err != nil {
		return ``, err
	}
	return string(hash), nil
}

// Verify compares a hashed password with its possible plaintext equivalent.
func (c *CryptoEngine) Verify(pass, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash),
		[]byte(pass+salt)); err != nil {
		return err
	}
	return nil
}
