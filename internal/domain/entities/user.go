// Package entities defines the core domain models for the application.
package entities

import (
	"errors"

	"github.com/aube/auth/internal/domain/valueobjects"
)

// User represents an application user account.
// Fields:
//   - ID: Database primary key
//   - Username: Unique identifier
//   - Email: Contact address
//   - Password: Hashed credentials (valueobjects.Password)
//
// Note: Excludes JSON tags to prevent accidental credential exposure
type User struct {
	ID       int64
	Username string
	Email    string
	Password *valueobjects.Password
}

// NewUser creates a validated User instance.
// id: Database ID (0 for new users)
// username: Unique handle
// email: Validated email address
// password: Hashed password object
// Returns: (*User, error) - validates all required fields
// Validation:
//   - Rejects empty username
//   - Rejects empty email
//   - Requires password object
func NewUser(id int64, username string, email string, password *valueobjects.Password) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if email == "" {
		return nil, errors.New("email cannot be nil")
	}

	if password == nil {
		return nil, errors.New("password cannot be nil")
	}

	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	}, nil
}

// SetPassword updates user credentials.
// newPassword: Pre-hashed password object
// Returns: error if password invalid
// Security: Enforces non-nil password requirement
func (u *User) SetPassword(newPassword *valueobjects.Password) error {
	if newPassword == nil {
		return errors.New("password cannot be nil")
	}

	u.Password = newPassword
	return nil
}

// PasswordMatches verifies credentials against stored hash.
// plainPassword: Candidate password
// Returns: bool indicating match
// Security: Uses constant-time comparison
func (u *User) PasswordMatches(plainPassword string) bool {
	return u.Password.Matches(plainPassword)
}

// GetHashedPassword retrieves the stored password hash.
// Returns: string representation of hash
// Security: Only returns the hashed value, never plaintext
func (u *User) GetHashedPassword() string {
	return u.Password.String()
}
