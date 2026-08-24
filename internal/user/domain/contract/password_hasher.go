// Password hasher contract.
//
// Intent: define password hashing and verification so credentials are never
// compared or stored in plain text.
// Objective: allow createuser to hash passwords on registration and login to
// verify a raw password against the stored hash.
package contract

// PasswordHasher defines hashing and verification of raw passwords.
type PasswordHasher interface {
	// Hash returns the hashed representation of a raw password.
	Hash(raw string) (string, error)
	// Verify reports whether the raw password matches the given hash.
	Verify(raw, hash string) (bool, error)
}
