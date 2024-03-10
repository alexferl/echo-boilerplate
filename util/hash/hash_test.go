package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	msg := "my message"
	key := "s3cret"
	h := NewHMAC([]byte(msg), []byte(key))
	b := ValidMAC([]byte(msg), []byte(h), []byte(key))

	assert.True(t, b)
}
