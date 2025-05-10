package model

import (
	"testing"
)

func TestUserValidate(t *testing.T) {
	t.Run("valid email and password", func(t *testing.T) {
		u := &User{
			Email:    "user@example.com",
			Password: "password123",
		}

		err := u.Validate()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		u := &User{
			Email:    "user@examplecom", // wrong domain part
			Password: "password123",
		}

		err := u.Validate()

		if err == nil {
			t.Error("expected error, got none")
		}
	})

	t.Run("empty password", func(t *testing.T) {
		u := &User{
			Email:    "user@example.com",
			Password: "", // empty password
		}

		err := u.Validate()

		if err == nil {
			t.Error("expected error, got none")
		}
	})
}
