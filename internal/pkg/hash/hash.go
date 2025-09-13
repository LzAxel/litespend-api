package hash

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordNotSame = errors.New("password is invalid")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrPasswordNotSame
	}

	return nil
}

func HashSessionKey(keyData string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(keyData), bcrypt.MinCost)
	return string(bytes), err
}
