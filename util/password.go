package util

import (
	"github.com/matthewhartstonge/argon2"
)

func HashPassword(password []byte) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded(password)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func VerifyPassword(password []byte, encoded []byte) error {
	_, err := argon2.VerifyEncoded(encoded, password)
	if err != nil {
		return err
	}

	return nil
}
