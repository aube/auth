package valueobjects

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	value string
}

func NewPassword(password string) (*Password, error) {
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}
	
	return &Password{value: password}, nil
}

func (p *Password) String() string {
	return p.value
}

func (p *Password) IsHashed() bool {
	// Простая проверка на то, хэширован ли пароль
	return len(p.value) == 60 && p.value[0] == '$'
}

func (p *Password) Hash() error {
	if p.IsHashed() {
		return nil
	}
	
	hashed, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	p.value = string(hashed)
	return nil
}

func (p *Password) Matches(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.value), []byte(plainPassword))
	return err == nil
}
