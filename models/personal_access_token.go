package models

import (
	"errors"
	"time"

	"github.com/rs/xid"

	"github.com/alexferl/echo-boilerplate/util/jwt"
	"github.com/alexferl/echo-boilerplate/util/password"
)

var ErrExpiresAtPast = errors.New("expires_at cannot be in the past")

type PersonalAccessToken struct {
	Id        string     `json:"id" bson:"id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
	IsRevoked bool       `json:"is_revoked" bson:"is_revoked"`
	Name      string     `json:"name" bson:"name"`
	Token     string     `json:"token" bson:"token"`
	UserId    string     `json:"user_id" bson:"user_id"`
}

func (pat *PersonalAccessToken) Encrypt() error {
	b, err := password.Hash([]byte(pat.Token))
	if err != nil {
		return err
	}

	pat.Token = b

	return nil
}

func (pat *PersonalAccessToken) Validate(s string) error {
	return password.Verify([]byte(pat.Token), []byte(s))
}

func (pat *PersonalAccessToken) Response() *PersonalAccessTokenResponse {
	return &PersonalAccessTokenResponse{
		Id:        pat.Id,
		CreatedAt: pat.CreatedAt,
		ExpiresAt: pat.ExpiresAt,
		Name:      pat.Name,
		Revoked:   pat.IsRevoked,
		UserId:    pat.UserId,
	}
}

type PersonalAccessTokens []PersonalAccessToken

type PersonalAccessTokensResponse struct {
	Tokens []PersonalAccessTokenResponse `json:"personal_access_tokens"`
}

func (pats PersonalAccessTokens) Response() *PersonalAccessTokensResponse {
	res := make([]PersonalAccessTokenResponse, 0)
	for _, pat := range pats {
		res = append(res, *pat.Response())
	}
	return &PersonalAccessTokensResponse{Tokens: res}
}

type PersonalAccessTokenResponse struct {
	Id        string     `json:"id" bson:"id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
	Name      string     `json:"name" bson:"name"`
	Revoked   bool       `json:"revoked" bson:"revoked"`
	UserId    string     `json:"user_id" bson:"user_id"`
}

func NewPersonalAccessToken(userId string, name string, expiresAt string) (*PersonalAccessToken, error) {
	t, err := time.Parse("2006-01-02", expiresAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if t.Before(now) {
		return nil, ErrExpiresAtPast
	}

	pat, err := jwt.GeneratePersonalToken(userId, t.Sub(now), nil)
	if err != nil {
		return nil, err
	}

	return &PersonalAccessToken{
		Id:        xid.New().String(),
		CreatedAt: &now,
		ExpiresAt: &t,
		Name:      name,
		Token:     string(pat),
		UserId:    userId,
	}, nil
}
