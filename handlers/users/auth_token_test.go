package users_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

type AuthTokenResponse struct {
	Exp   time.Time `json:"exp"`
	Iat   time.Time `json:"iat"`
	Iss   string    `json:"iss"`
	Nbf   time.Time `json:"nbf"`
	Roles []string  `json:"roles"`
	Sub   string    `json:"sub"`
	Type  string    `json:"type"`
}

func TestHandler_AuthToken_200(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := util.ParseToken(access)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	var result AuthTokenResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)

	roles, _ := token.Get("roles")
	typ, _ := token.Get("type")

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, token.Expiration(), result.Exp)
	assert.Equal(t, token.IssuedAt(), result.Iat)
	assert.Equal(t, token.Issuer(), result.Iss)
	assert.Equal(t, token.NotBefore(), result.Nbf)
	assert.ElementsMatch(t, roles, user.Roles)
	assert.Equal(t, token.Subject(), result.Sub)
	assert.Equal(t, typ, result.Type)
}

func TestHandler_AuthToken_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_AuthToken_200_Cookie(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := util.ParseToken(access)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/auth/token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(util.NewAccessTokenCookie(access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	var result AuthTokenResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)

	roles, _ := token.Get("roles")
	typ, _ := token.Get("type")

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, token.Expiration(), result.Exp)
	assert.Equal(t, token.IssuedAt(), result.Iat)
	assert.Equal(t, token.Issuer(), result.Iss)
	assert.Equal(t, token.NotBefore(), result.Nbf)
	assert.ElementsMatch(t, roles, user.Roles)
	assert.Equal(t, token.Subject(), result.Sub)
	assert.Equal(t, typ, result.Type)
}
