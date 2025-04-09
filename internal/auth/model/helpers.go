package model

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ParseCredentials(requestBody []byte) (User, error) {
	var user User

	err := json.Unmarshal(requestBody, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
