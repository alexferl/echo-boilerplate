package util

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/minio/sha256-simd"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
)

var ErrTokenParse = errors.New("failed parsing token")

type TokenType int64

const (
	AccessToken TokenType = iota + 1
	RefreshToken
)

func (t TokenType) String() string {
	return [...]string{"access", "refresh"}[t-1]
}

func GenerateTokens(sub string, claims map[string]any) ([]byte, []byte, error) {
	access, err := GenerateAccessToken(sub, claims)
	if err != nil {
		return nil, nil, err
	}

	refresh, err := GenerateRefreshToken(sub)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}

func GenerateAccessToken(sub string, claims map[string]any) ([]byte, error) {
	return generateToken(AccessToken, sub, claims)
}

func GenerateRefreshToken(sub string) ([]byte, error) {
	return generateToken(RefreshToken, sub, map[string]any{})
}

func generateToken(typ TokenType, sub string, claims map[string]any) ([]byte, error) {
	key, err := LoadPrivateKey()
	if err != nil {
		return nil, err
	}

	var expiry time.Duration
	switch typ {
	case AccessToken:
		expiry = viper.GetDuration(config.JWTAccessTokenExpiry)
	case RefreshToken:
		expiry = viper.GetDuration(config.JWTRefreshTokenExpiry)
	default:
		return nil, fmt.Errorf("invalid token type")
	}

	builder := jwt.NewBuilder().
		Subject(sub).
		Issuer(viper.GetString(config.JWTIssuer)).
		IssuedAt(time.Now()).
		NotBefore(time.Now()).
		Expiration(time.Now().Add(expiry)).
		Claim("type", typ.String())

	for k, v := range claims {
		builder.Claim(k, v)
	}

	token, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build %s token: %v\n", typ.String(), err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256, key))
	if err != nil {
		return nil, fmt.Errorf("failed to sign %s token: %v\n", typ.String(), err)
	}

	return signed, nil
}

func ParseToken(encodedToken string) (jwt.Token, error) {
	key, err := LoadPrivateKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse([]byte(encodedToken), jwt.WithValidate(true), jwt.WithKey(jwa.RS256, key))
	if err != nil {
		if err == jwt.ErrTokenExpired() {
			return nil, err
		}

		if err == jwt.ErrInvalidIssuedAt() {
			return nil, err
		}

		if err == jwt.ErrTokenNotYetValid() {
			return nil, err
		}

		return nil, ErrTokenParse
	}

	return token, nil
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	f, err := os.Open(viper.GetString(config.JWTPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to open private key: %v", err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block: %v", err)
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func HashToken(token []byte) string {
	h := sha256.New()
	h.Write(token)
	b := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(b)
}
