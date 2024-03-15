package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
)

type PersonalAccessTokenTestSuite struct {
	suite.Suite
	mapper *services.MockPersonalAccessTokenMapper
	svc    *services.PersonalAccessToken
	user   *models.User
}

func (s *PersonalAccessTokenTestSuite) SetupTest() {
	s.mapper = services.NewMockPersonalAccessTokenMapper(s.T())
	s.svc = services.NewPersonalAccessToken(s.mapper)
	user := models.NewUser("test@email.com", "test")
	user.Id = "100"
	s.user = user
}

func TestPersonalAccessToken(t *testing.T) {
	suite.Run(t, new(PersonalAccessTokenTestSuite))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessToken_Create() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.user.Id, name, expiresAt)
	s.Assert().NoError(err)

	s.mapper.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.Create(context.Background(), m)
	s.Assert().NoError(err)
	s.Assert().Equal(name, pat.Name)
	s.Assert().Equal(expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_Read() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.user.Id, name, expiresAt)
	s.Assert().NoError(err)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.Read(context.Background(), s.user.Id, id)
	s.Assert().NoError(err)
	s.Assert().Equal(id, pat.Id)
	s.Assert().Equal(name, pat.Name)
	s.Assert().Equal(expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_Read_Err() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.user.Id, name, expiresAt)
	s.Assert().NoError(err)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	_, err = s.svc.Read(context.Background(), s.user.Id, id)
	s.Assert().Error(err)
	var se *services.Error
	s.Assert().ErrorAs(err, &se)
	if errors.As(err, &se) {
		s.Assert().Equal(services.NotExist, se.Kind)
	}
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_Revoke() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.user.Id, name, expiresAt)
	s.Assert().NoError(err)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(m, nil)

	err = s.svc.Revoke(context.Background(), m)
	s.Assert().NoError(err)

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.Read(context.Background(), s.user.Id, id)
	s.Assert().NoError(err)
	s.Assert().True(pat.IsRevoked)
	s.Assert().Equal(id, pat.Id)
	s.Assert().Equal(name, pat.Name)
	s.Assert().Equal(expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_Find() {
	s.mapper.EXPECT().
		Find(mock.Anything, mock.Anything).
		Return(models.PersonalAccessTokens{}, nil)

	pats, err := s.svc.Find(context.Background(), "123")
	s.Assert().NoError(err)
	s.Assert().Equal(models.PersonalAccessTokens{}, pats)
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_FindOne() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.user.Id, name, expiresAt)
	s.Assert().NoError(err)
	id := "123"
	userId := "456"
	m.Id = id
	m.UserId = userId

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	pat, err := s.svc.FindOne(context.Background(), userId, name)
	s.Assert().NoError(err)
	s.Assert().Equal(id, pat.Id)
	s.Assert().Equal(userId, pat.UserId)
	s.Assert().Equal(name, pat.Name)
	s.Assert().Equal(expiresAt, pat.ExpiresAt.Format("2006-01-02"))
}

func (s *PersonalAccessTokenTestSuite) TestPersonalAccessTokenTestSuite_FindOne_Err() {
	name := "my_token"
	expiresAt := time.Now().Add((7 * 24) * time.Hour).Format("2006-01-02")
	m, err := models.NewPersonalAccessToken(s.user.Id, name, expiresAt)
	s.Assert().NoError(err)
	id := "123"
	userId := "456"
	m.Id = id
	m.UserId = userId

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	_, err = s.svc.FindOne(context.Background(), id, "")
	s.Assert().Error(err)
	var se *services.Error
	s.Assert().ErrorAs(err, &se)
	if errors.As(err, &se) {
		s.Assert().Equal(services.NotExist, se.Kind)
	}
}
