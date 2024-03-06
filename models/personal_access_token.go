package models

import (
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/xid"

	"github.com/alexferl/echo-boilerplate/util"
)

var ErrExpiresAtPast = errors.New("expires_at cannot be in the past")

type PersonalAccessToken struct {
	Id        string     `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Revoked   bool       `json:"revoked" bson:"revoked"`
	UserId    string     `json:"user_id" bson:"user_id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
	Token     string     `json:"token" bson:"token"`
}

func (pat *PersonalAccessToken) Encrypt() error {
	b, err := util.HashPassword([]byte(pat.Token))
	if err != nil {
		return err
	}

	pat.Token = b

	return nil
}

func (pat *PersonalAccessToken) Validate(s string) error {
	return util.VerifyPassword([]byte(pat.Token), []byte(s))
}

func (pat *PersonalAccessToken) Response() *PersonalAccessTokenResponse {
	return &PersonalAccessTokenResponse{
		Id:        pat.Id,
		Name:      pat.Name,
		Revoked:   pat.Revoked,
		UserId:    pat.UserId,
		CreatedAt: pat.CreatedAt,
		ExpiresAt: pat.ExpiresAt,
	}
}

type PersonalAccessTokens []PersonalAccessToken

func (pats PersonalAccessTokens) Response() []*PersonalAccessTokenResponse {
	res := make([]*PersonalAccessTokenResponse, 0)
	for _, pat := range pats {
		res = append(res, pat.Response())
	}
	return res
}

type PersonalAccessTokenResponse struct {
	Id        string     `json:"id" bson:"id"`
	Name      string     `json:"name" bson:"name"`
	Revoked   bool       `json:"revoked" bson:"revoked"`
	UserId    string     `json:"user_id" bson:"user_id"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	ExpiresAt *time.Time `json:"expires_at" bson:"expires_at"`
}

func NewPersonalAccessToken(token jwt.Token, name string, expiresAt string) (*PersonalAccessToken, error) {
	t, err := time.Parse("2006-01-02", expiresAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	if t.Before(now) {
		return nil, ErrExpiresAtPast
	}

	roles := util.GetRoles(token)
	pat, err := util.GeneratePersonalToken(token.Subject(), t.Sub(now), map[string]any{"roles": roles})
	if err != nil {
		return nil, err
	}

	return &PersonalAccessToken{
		Id:        xid.New().String(),
		Name:      name,
		UserId:    token.Subject(),
		Token:     string(pat),
		CreatedAt: &now,
		ExpiresAt: &t,
	}, nil
}
