package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"slices"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	jwx "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/alexferl/echo-boilerplate/config"
)

type TokenType int64

const (
	AccessToken TokenType = iota + 1
	RefreshToken
	PersonalToken
)

func (t TokenType) String() string {
	return [...]string{"access", "refresh", "personal"}[t-1]
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
	expiry := viper.GetDuration(config.JWTAccessTokenExpiry)
	return generateToken(AccessToken, expiry, sub, claims)
}

func GenerateRefreshToken(sub string) ([]byte, error) {
	expiry := viper.GetDuration(config.JWTRefreshTokenExpiry)
	return generateToken(RefreshToken, expiry, sub, map[string]any{})
}

func GeneratePersonalToken(sub string, expiry time.Duration, claims map[string]any) ([]byte, error) {
	return generateToken(PersonalToken, expiry, sub, claims)
}

func generateToken(typ TokenType, expiry time.Duration, sub string, claims map[string]any) ([]byte, error) {
	key, err := LoadPrivateKey()
	if err != nil {
		return nil, err
	}

	builder := jwx.NewBuilder().
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

	signed, err := jwx.Sign(token, jwx.WithKey(jwa.RS256, key))
	if err != nil {
		return nil, fmt.Errorf("failed to sign %s token: %v\n", typ.String(), err)
	}

	return signed, nil
}

func ParseEncoded(encodedToken []byte) (jwx.Token, error) {
	key, err := LoadPrivateKey()
	if err != nil {
		return nil, err
	}

	token, err := jwx.Parse(encodedToken, jwx.WithValidate(true), jwx.WithKey(jwa.RS256, key))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetRoles(token jwx.Token) []string {
	val, _ := token.Get("roles")
	roles := val.([]interface{})

	var res []string
	for _, role := range roles {
		res = append(res, role.(string))
	}

	return res
}

func HasRoles(token jwx.Token, roles ...string) bool {
	if val, ok := token.Get("roles"); !ok {
		log.Error().Msgf("failed getting roles for token: %s", token.Subject())
		return false
	} else {
		jwtRoles, ok := val.([]interface{})
		if !ok {
			log.Error().Msgf("failed converting roles to slice for token: %s", token.Subject())
			return false
		}

		for _, r := range jwtRoles {
			if slices.Contains(roles, r.(string)) {
				return true
			}
		}
	}

	return false
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

	var key *rsa.PrivateKey
	switch block.Type {
	case "RSA PRIVATE KEY": // PKCS#1
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	case "PRIVATE KEY": // PKCS#8
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		key = privateKey.(*rsa.PrivateKey)
	}

	return key, nil
}
