// Package valueobjects contains domain primitives that enforce business rules.
package valueobjects

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// Email represents a validated email address value object.
// Implements proper encapsulation and validation of email addresses.
type Email struct {
	value string
}

// NewEmail creates a validated Email instance.
// value: Candidate email address
// Returns: (Email, error)
// Validation:
//   - Requires 3-254 characters
//   - Valid format per RFC 5322
//   - Normalizes to lowercase
//   - Trims whitespace
func NewEmail(value string) (Email, error) {
	if !isValidEmail(value) {
		return Email{}, errors.New("invalid email address")
	}

	return Email{
		value: strings.ToLower(strings.TrimSpace(value)),
	}, nil
}

func isValidEmail(email string) bool {
	// Basic length check
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	// Use Go's built-in email parser
	_, err := mail.ParseAddress(email)
	return err == nil
}

// String returns the normalized email string.
// Implements fmt.Stringer interface.
func (e Email) String() string {
	return e.value
}

// Equals compares two Email objects for value equality.
// other: Email to compare
// Returns: bool indicating equality
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// LocalPart extracts the username portion (before @).
// Returns: string - empty if malformed
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// Domain extracts the domain portion (after @).
// Returns: string - empty if malformed
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// MarshalText implements encoding.TextMarshaler for serialization.
// Returns: ([]byte, error) - UTF-8 encoded email
func (e Email) MarshalText() ([]byte, error) {
	return []byte(e.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for deserialization.
// text: Input byte slice
// Returns: error if validation fails
func (e *Email) UnmarshalText(text []byte) error {
	tmp, err := NewEmail(string(text))
	if err != nil {
		return err
	}
	e.value = tmp.value
	return nil
}

// Value returns the underlying string value.
// Returns: string - normalized email
func (e Email) Value() string {
	return e.value
}

// IsZero checks for uninitialized state.
// Returns: bool - true if empty
func (e Email) IsZero() bool {
	return e.value == ""
}

// Format implements custom formatting for fmt package.
// Supports:
//   - %v: Basic string output
//   - %+v: Debug format with type information
func (e Email) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "Email{value: %q}", e.value)
		} else {
			fmt.Fprint(f, e.value)
		}
	default:
		fmt.Fprint(f, e.value)
	}
}
