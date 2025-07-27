package entities_test

import (
	"errors"
	"testing"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_Success(t *testing.T) {
	password, _ := valueobjects.NewPassword("password123")
	user, err := entities.NewUser(1, "testuser", "test@example.com", password)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, password, user.Password)
}

func TestNewUser_ValidationErrors(t *testing.T) {
	password, _ := valueobjects.NewPassword("password123")

	tests := []struct {
		name     string
		username string
		email    string
		password *valueobjects.Password
		expected error
	}{
		{"empty username", "", "test@example.com", password, errors.New("username cannot be empty")},
		{"empty email", "testuser", "", password, errors.New("email cannot be nil")},
		{"nil password", "testuser", "test@example.com", nil, errors.New("password cannot be nil")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := entities.NewUser(1, tt.username, tt.email, tt.password)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expected.Error())
		})
	}
}

func TestUserPasswordMethods(t *testing.T) {
	// Создаем тестовый пароль
	password, err := valueobjects.NewPassword("password123")
	require.NoError(t, err)

	// Хэшируем пароль
	err = password.Hash()
	require.NoError(t, err)

	user, err := entities.NewUser(1, "testuser", "test@example.com", password)
	require.NoError(t, err)

	t.Run("PasswordMatches", func(t *testing.T) {
		assert.True(t, user.PasswordMatches("password123"))
		assert.False(t, user.PasswordMatches("wrongpassword"))
	})

	t.Run("GetHashedPassword", func(t *testing.T) {
		hashed := user.GetHashedPassword()
		assert.NotEmpty(t, hashed)
		assert.NotEqual(t, "password123", hashed)
	})

	t.Run("SetPassword", func(t *testing.T) {
		newPassword, err := valueobjects.NewPassword("newpassword123")
		require.NoError(t, err)
		err = newPassword.Hash()
		require.NoError(t, err)

		err = user.SetPassword(newPassword)
		assert.NoError(t, err)
		assert.True(t, user.PasswordMatches("newpassword123"))

		err = user.SetPassword(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be nil")
	})
}
