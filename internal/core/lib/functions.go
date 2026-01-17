package lib

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCifrarPassword = errors.New("error al cifrar contraceña")
)

func EncryptPassword(pasword string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pasword), bcrypt.DefaultCost)

	if err != nil {
		return "", ErrCifrarPassword
	}

	return string(hashedPassword), nil
}
