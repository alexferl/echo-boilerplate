package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexferl/echo-boilerplate/data"

	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

func TestHandler_CreatePersonalAccessToken_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := util.ParseToken(access)
	assert.NoError(t, err)

	payload := &users.CreatePATRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	newPAT, err := users.NewPersonalAccessToken(token, payload.Name, payload.ExpiresAt)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/user/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			nil,
		).
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"FindOneAndUpdate",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			newPAT,
			nil,
		)

	s.ServeHTTP(resp, req)

	var result users.PersonalAccessToken
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_CreatePersonalAccessToken_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodPost, "/user/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_CreatePersonalAccessToken_409(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := util.ParseToken(access)
	assert.NoError(t, err)

	payload := &users.CreatePATRequest{
		Name:      "My Token",
		ExpiresAt: time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02"),
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	newPAT, err := users.NewPersonalAccessToken(token, payload.Name, payload.ExpiresAt)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/user/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			newPAT,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)
}

func TestHandler_CreatePersonalAccessToken_422(t *testing.T) {
	_, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	payload := &users.CreatePATRequest{
		Name:      "",
		ExpiresAt: "invalid",
	}
	b, err := json.Marshal(payload)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/user/personal_access_tokens", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func createTokens(t *testing.T, token jwt.Token, num int) []*users.PATWithoutToken {
	var result []*users.PATWithoutToken

	for i := 1; i <= num; i++ {
		pat, err := users.NewPersonalAccessToken(
			token,
			fmt.Sprintf("my_token%d", i),
			time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
		)
		assert.NoError(t, err)
		resp := pat.MakeResponse()
		result = append(result, resp)
	}

	return result
}

func TestHandler_ListPersonalAccessTokens_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := util.ParseToken(access)
	assert.NoError(t, err)

	tokens := createTokens(t, token, 10)

	req := httptest.NewRequest(http.MethodGet, "/user/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"Find",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			tokens,
			nil,
		)

	s.ServeHTTP(resp, req)

	var result users.ListPATResponse
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, 10, len(result.Tokens))
}

func TestHandler_ListPersonalAccessTokens_401(t *testing.T) {
	_, s := getMapperAndServer(t)

	req := httptest.NewRequest(http.MethodGet, "/user/personal_access_tokens", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHandler_GetPersonalAccessToken_200(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := util.ParseToken(access)
	assert.NoError(t, err)

	newPAT, err := users.NewPersonalAccessToken(
		token,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)
	assert.NoError(t, err)
	pat := newPAT.MakeResponse()

	req := httptest.NewRequest(http.MethodGet, "/user/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			pat,
			nil,
		)

	s.ServeHTTP(resp, req)

	var result users.PATWithoutToken
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandler_GetPersonalAccessToken_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/user/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			data.ErrNoDocuments,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestHandler_RevokePersonalAccessToken_204(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	token, err := util.ParseToken(access)
	assert.NoError(t, err)

	newPAT, err := users.NewPersonalAccessToken(
		token,
		fmt.Sprintf("my_token"),
		time.Now().Add((7*24)*time.Hour).Format("2006-01-02"),
	)
	assert.NoError(t, err)
	pat := newPAT.MakeResponse()

	req := httptest.NewRequest(http.MethodDelete, "/user/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			pat,
			nil,
		).
		On(
			"UpdateOneById",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			nil,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func TestHandler_RevokePersonalAccessToken_404(t *testing.T) {
	mapper, s := getMapperAndServer(t)

	user := users.NewUser("test@example.com", "test")
	access, _, err := user.Login()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/user/personal_access_tokens/id", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
	resp := httptest.NewRecorder()

	mapper.Mock.
		On(
			"WithCollection",
			mock.Anything,
		).
		Return(
			mapper,
		).
		On(
			"FindOne",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).
		Return(
			nil,
			data.ErrNoDocuments,
		)

	s.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
