package model

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestEncryptString(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		s := "password123"

		encrypted, err := encryptString(s)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(s))
		if err != nil {
			t.Error("expected valid encrypted string")
		}
	})
}
