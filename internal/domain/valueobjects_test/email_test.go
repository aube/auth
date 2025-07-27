package valueobjects_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aube/auth/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmail_ValidEmails(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"simple", "test@example.com"},
		{"with dots", "first.last@example.com"},
		{"with plus", "test+filter@example.com"},
		{"subdomain", "test@sub.example.com"},
		{"uppercase", "TEST@example.com"},
		{"with spaces", "  test@example.com  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := valueobjects.NewEmail(tt.email)
			require.NoError(t, err)
			assert.Equal(t, strings.ToLower(strings.TrimSpace(tt.email)), email.Value())
		})
	}
}

func TestNewEmail_InvalidEmails(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"empty", ""},
		{"no at", "invalid.com"},
		{"no domain", "test@"},
		{"no local", "@example.com"},
		{"too long", strings.Repeat("a", 255) + "@example.com"},
		{"invalid chars", "test@exa mple.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := valueobjects.NewEmail(tt.email)
			assert.Error(t, err)
		})
	}
}

func TestEmail_Methods(t *testing.T) {
	email, err := valueobjects.NewEmail("Test.User+filter@Example.COM")
	require.NoError(t, err)

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "test.user+filter@example.com", email.String())
	})

	t.Run("Value", func(t *testing.T) {
		assert.Equal(t, "test.user+filter@example.com", email.Value())
	})

	t.Run("Equals", func(t *testing.T) {
		other, _ := valueobjects.NewEmail("test.user+filter@example.com")
		assert.True(t, email.Equals(other))
	})

	t.Run("LocalPart", func(t *testing.T) {
		assert.Equal(t, "test.user+filter", email.LocalPart())
	})

	t.Run("Domain", func(t *testing.T) {
		assert.Equal(t, "example.com", email.Domain())
	})

	t.Run("IsZero", func(t *testing.T) {
		assert.False(t, email.IsZero())
		zero := valueobjects.Email{}
		assert.True(t, zero.IsZero())
	})
}

func TestEmail_Marshaling(t *testing.T) {
	email, _ := valueobjects.NewEmail("test@example.com")

	t.Run("MarshalText", func(t *testing.T) {
		data, err := email.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", string(data))
	})

	t.Run("UnmarshalText", func(t *testing.T) {
		var e valueobjects.Email
		err := e.UnmarshalText([]byte("test@example.com"))
		require.NoError(t, err)
		assert.Equal(t, "test@example.com", e.Value())
	})

	t.Run("UnmarshalText invalid", func(t *testing.T) {
		var e valueobjects.Email
		err := e.UnmarshalText([]byte("invalid"))
		assert.Error(t, err)
	})
}

func TestEmail_Format(t *testing.T) {
	email, _ := valueobjects.NewEmail("test@example.com")

	t.Run("default format", func(t *testing.T) {
		assert.Equal(t, "test@example.com", email.String())
	})

	t.Run("verbose format", func(t *testing.T) {
		assert.Equal(t, `Email{value: "test@example.com"}`, fmt.Sprintf("%+v", email))
	})
}
