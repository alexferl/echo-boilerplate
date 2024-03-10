package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func NewHMAC(message []byte, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

// ValidMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidMAC(message []byte, messageMAC []byte, key []byte) bool {
	expectedMAC := NewHMAC(message, key)
	return hmac.Equal(messageMAC, []byte(expectedMAC))
}
