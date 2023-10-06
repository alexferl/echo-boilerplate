package util

import (
	"github.com/matthewhartstonge/argon2"
)

func HashPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func VerifyPassword(password string, encoded string) error {
	_, err := argon2.VerifyEncoded([]byte(encoded), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
