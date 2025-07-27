package valueobjects_test

import (
	"testing"

	"github.com/aube/auth/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"valid", "password123", true},
		{"too short", "short", false},
		{"empty", "", false},
		{"exact min length", "12345678", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := valueobjects.NewPassword(tt.password)
			if tt.valid {
				require.NoError(t, err)
				assert.Equal(t, tt.password, p.String())
				assert.False(t, p.IsHashed())
			} else {
				assert.Error(t, err)
				assert.Nil(t, p)
			}
		})
	}
}

func TestPassword_Hash(t *testing.T) {
	p, _ := valueobjects.NewPassword("password123")

	t.Run("hash new password", func(t *testing.T) {
		err := p.Hash()
		require.NoError(t, err)
		assert.True(t, p.IsHashed())
		assert.True(t, len(p.String()) == 60)
	})

	t.Run("hash already hashed", func(t *testing.T) {
		original := p.String()
		err := p.Hash()
		assert.NoError(t, err)
		assert.Equal(t, original, p.String())
	})
}

func TestPassword_Matches(t *testing.T) {
	p, _ := valueobjects.NewPassword("password123")
	require.NoError(t, p.Hash())

	t.Run("correct password", func(t *testing.T) {
		assert.True(t, p.Matches("password123"))
	})

	t.Run("incorrect password", func(t *testing.T) {
		assert.False(t, p.Matches("wrongpassword"))
	})

	t.Run("not hashed", func(t *testing.T) {
		p2, _ := valueobjects.NewPassword("password123")
		assert.False(t, p2.Matches("password123"))
	})
}
