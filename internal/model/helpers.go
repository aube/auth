package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func generateRandomString(n int, includeSpecial bool) (string, error) {
	const (
		letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	)

	charset := letterBytes
	if includeSpecial {
		charset += specialChars
	}

	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		ret[i] = charset[num.Int64()]
	}

	return string(ret), nil
}

func makeHash() (string, error) {
	input, err := generateRandomString(32, true)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}

func ParseCredentials(requestBody []byte) (User, error) {
	var user User

	err := json.Unmarshal(requestBody, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
