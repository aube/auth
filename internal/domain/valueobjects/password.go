// Package valueobjects contains domain primitives for security-sensitive values.
package valueobjects

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Password represents a secure password value object.
// Encapsulates hashing and verification logic.
type Password struct {
	value string
}

// NewPassword creates a validated Password instance.
// password: Plaintext password
// Returns: (*Password, error)
// Validation:
//   - Minimum 8 characters
//   - Returns error for weak passwords
func NewPassword(password string) (*Password, error) {
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	return &Password{value: password}, nil
}

// String returns the current value (hashed or plaintext).
// Note: Avoid logging or exposing this value.
// Implements fmt.Stringer interface.
func (p *Password) String() string {
	return p.value
}

// IsHashed checks if the password is already hashed.
// Heuristic:
//   - 60 characters
//   - Starts with bcrypt prefix ($2a$)
//
// Returns: bool
func (p *Password) IsHashed() bool {
	// Простая проверка на то, хэширован ли пароль
	return len(p.value) == 60 && p.value[0] == '$'
}

// Hash securely hashes the password using bcrypt.
// Idempotent: skips if already hashed.
// Returns: error on hashing failure
// Security:
//   - Uses bcrypt.DefaultCost (currently 10)
//   - Salt is automatically generated
func (p *Password) Hash() error {
	if p.IsHashed() {
		return nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.value = string(hashed)
	return nil
}

// Matches verifies a plaintext password against the stored value.
// plainPassword: Candidate password
// Returns: bool indicating match
// Security:
//   - Uses constant-time comparison
//   - Handles both hashed and plaintext states
func (p *Password) Matches(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.value), []byte(plainPassword))
	return err == nil
}
