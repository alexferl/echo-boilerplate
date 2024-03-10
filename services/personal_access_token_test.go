package services_test

import (
	"context"
	"testing"
	"time"

	jwx "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
	"github.com/alexferl/echo-boilerplate/util/jwt"
)

type PersonalAccessTokenTestSuite struct {
	suite.Suite
	mapper *services.MockPersonalAccessTokenMapper
	svc    *services.PersonalAccessToken
	token  jwx.Token
}

func (s *PersonalAccessTokenTestSuite) SetupTest() {
	s.mapper = services.NewMockPersonalAccessTokenMapper(s.T())
	s.svc = services.NewPersonalAccessToken(s.mapper)
	user := models.NewUser("test@email.com", "test")
	access, _, _ := user.Login()
	token, _ := jwt.ParseEncoded(access)
	s.token = token
}

func TestPersonalAccessToken(t *testing.T) {
	suite.Run(t, new(PersonalAccessTokenTestSuite))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessToken_Create() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.token, name, expiresAt)
	assert.NoError(s.T(), err)

	s.mapper.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.Create(context.Background(), m)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), name, pat.Name)
	assert.Equal(s.T(), expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_Read() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.token, name, expiresAt)
	assert.NoError(s.T(), err)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.Read(context.Background(), id)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), id, pat.Id)
	assert.Equal(s.T(), name, pat.Name)
	assert.Equal(s.T(), expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_Revoke() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.token, name, expiresAt)
	assert.NoError(s.T(), err)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(m, nil)

	err = s.svc.Revoke(context.Background(), m)
	assert.NoError(s.T(), err)

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.Read(context.Background(), id)
	assert.NoError(s.T(), err)

	assert.True(s.T(), pat.IsRevoked)
	assert.Equal(s.T(), id, pat.Id)
	assert.Equal(s.T(), name, pat.Name)
	assert.Equal(s.T(), expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_Find() {
	s.mapper.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(models.PersonalAccessTokens{}, nil)

	pats, err := s.svc.Find(context.Background(), "123")
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), models.PersonalAccessTokens{}, pats)
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_FindOne() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.token, name, expiresAt)
	assert.NoError(s.T(), err)
	id := "123"
	userId := "456"
	m.Id = id
	m.UserId = userId

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.FindOne(context.Background(), userId, name)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), id, pat.Id)
	assert.Equal(s.T(), userId, pat.UserId)
	assert.Equal(s.T(), name, pat.Name)
	assert.Equal(s.T(), expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}
