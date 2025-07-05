package entities

import (
	"errors"

	"github.com/aube/auth/internal/domain/valueobjects"
)

type User struct {
	ID       int64
	Username string
	password *valueobjects.Password
}

func NewUser(id int64, username string, password *valueobjects.Password) (*User, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	if password == nil {
		return nil, errors.New("password cannot be nil")
	}

	return &User{
		ID:       id,
		Username: username,
		password: password,
	}, nil
}

func (u *User) SetPassword(newPassword *valueobjects.Password) error {
	if newPassword == nil {
		return errors.New("password cannot be nil")
	}

	u.password = newPassword
	return nil
}

func (u *User) PasswordMatches(plainPassword string) bool {
	return u.password.Matches(plainPassword)
}

func (u *User) GetHashedPassword() string {
	return u.password.String()
}
