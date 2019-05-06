package util

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const charset = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// HashPassword returns a bcrypt hash
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks that the password and hash match (password is valid)
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomString returns a random string of length n
func RandomString(length int) string {
	return stringWithCharset(length, charset)
}
