// Bcrypt password hasher adapter.
//
// Intent: implement the domain's contract.PasswordHasher using the bcrypt algorithm.
// Objective: hash raw passwords for persistence and verify a raw password against a stored
// bcrypt hash, without exposing bcrypt-specific details to the domain or application layers.
package database

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptAdapter implements contract.PasswordHasher using bcrypt.
type BcryptAdapter struct{}

// NewBcryptAdapter creates a BcryptAdapter.
//
// Intent: construct a ready-to-use bcrypt-based password hasher.
// Returns: a *BcryptAdapter satisfying contract.PasswordHasher.
func NewBcryptAdapter() *BcryptAdapter {
	return &BcryptAdapter{}
}

// Hash returns the bcrypt hash of a raw password.
//
// Intent: derive a salted bcrypt hash suitable for persistence, using bcrypt's default cost
// factor.
// Parameters:
//   - raw: the plain-text password to hash.
//
// Returns: the bcrypt hash as a string, or an error if bcrypt fails to generate it (e.g. the
// raw password exceeds bcrypt's 72-byte input limit).
func (a *BcryptAdapter) Hash(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify reports whether the raw password matches the given bcrypt hash.
//
// Intent: authenticate a raw password against its previously stored bcrypt hash without ever
// comparing plain text.
// Parameters:
//   - raw: the plain-text password supplied by the caller.
//   - hash: the bcrypt hash previously produced by Hash.
//
// Returns: true if raw matches hash; false (with a nil error) when they simply don't match; a
// non-nil error only for unexpected bcrypt failures (e.g. a malformed hash).
func (a *BcryptAdapter) Verify(raw, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
