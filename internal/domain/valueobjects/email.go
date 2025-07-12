package valueobjects

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

// Email represents a validated email address value object
type Email struct {
	value string
}

// NewEmail creates a new Email value object after validation
func NewEmail(value string) (Email, error) {
	if !isValidEmail(value) {
		return Email{}, errors.New("invalid email address")
	}

	return Email{
		value: strings.ToLower(strings.TrimSpace(value)),
	}, nil
}

// isValidEmail checks if the email address is valid
func isValidEmail(email string) bool {
	// Basic length check
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Use Go's built-in email parser
	_, err := mail.ParseAddress(email)
	return err == nil
}

// String returns the string representation of the email
func (e Email) String() string {
	return e.value
}

// Equals checks if two Email objects are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// LocalPart returns the local part of the email (before @)
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// Domain returns the domain part of the email (after @)
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// MarshalText implements the encoding.TextMarshaler interface
func (e Email) MarshalText() ([]byte, error) {
	return []byte(e.value), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface
func (e *Email) UnmarshalText(text []byte) error {
	tmp, err := NewEmail(string(text))
	if err != nil {
		return err
	}
	e.value = tmp.value
	return nil
}

// Value returns the underlying string value
func (e Email) Value() string {
	return e.value
}

// IsZero checks if the email is zero value
func (e Email) IsZero() bool {
	return e.value == ""
}

// Format implements fmt.Formatter interface
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
