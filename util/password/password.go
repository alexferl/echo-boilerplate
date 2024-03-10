package password

import (
	"errors"

	"github.com/matthewhartstonge/argon2"
)

func Hash(password []byte) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded(password)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func Verify(password []byte, encoded []byte) error {
	b, err := argon2.VerifyEncoded(encoded, password)
	if err != nil {
		return err
	}

	if !b {
		return errors.New("mismatch")
	}

	return nil
}
