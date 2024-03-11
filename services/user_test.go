package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
)

type UserTestSuite struct {
	suite.Suite
	mapper *services.MockUserMapper
	svc    *services.User
}

func (s *UserTestSuite) SetupTest() {
	s.mapper = services.NewMockUserMapper(s.T())
	s.svc = services.NewUser(s.mapper)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestUser_Create() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Create(id)

	s.mapper.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(m, nil)

	user, err := s.svc.Create(context.Background(), m)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user.CreatedBy)
	assert.Equal(s.T(), email, user.Email)
	assert.Equal(s.T(), username, user.Username)
}

func (s *UserTestSuite) TestUser_Create_Err() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Create(id)

	s.mapper.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil, &mongo.WriteError{Code: 11000})

	_, err := s.svc.Create(context.Background(), m)
	assert.Error(s.T(), err)
	var se *services.Error
	assert.ErrorAs(s.T(), err, &se)
	if errors.As(err, &se) {
		assert.Equal(s.T(), services.Exist, se.Kind)
	}
}

func (s *UserTestSuite) TestUser_Read() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	user, err := s.svc.Read(context.Background(), id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), id, user.Id)
	assert.Equal(s.T(), email, user.Email)
	assert.Equal(s.T(), username, user.Username)
}

func (s *UserTestSuite) TestUser_Read_Err() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	_, err := s.svc.Read(context.Background(), id)
	assert.Error(s.T(), err)
	var se *services.Error
	assert.ErrorAs(s.T(), err, &se)
	if errors.As(err, &se) {
		assert.Equal(s.T(), services.NotExist, se.Kind)
	}
}

func (s *UserTestSuite) TestUser_Update() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Update(id)

	s.mapper.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(m, nil)

	task, err := s.svc.Update(context.Background(), id, m)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), task.UpdatedBy)
}

func (s *UserTestSuite) TestUser_Delete() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Delete(id)

	s.mapper.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(m, nil)

	err := s.svc.Delete(context.Background(), id, m)
	assert.NoError(s.T(), err)

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	_, err = s.svc.Read(context.Background(), id)
	assert.Error(s.T(), err)
	var se *services.Error
	assert.ErrorAs(s.T(), err, &se)
	if errors.As(err, &se) {
		assert.Equal(s.T(), services.Deleted, se.Kind)
	}
}

func (s *UserTestSuite) TestUser_Find() {
	s.mapper.EXPECT().
		Find(mock.Anything, mock.Anything, 1, 0).
		Return(1, models.Users{}, nil)

	count, tasks, err := s.svc.Find(context.Background(), &models.UserSearchParams{
		Limit: 1,
		Skip:  0,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), count)
	assert.Equal(s.T(), models.Users{}, tasks)
}

func (s *UserTestSuite) TestUser_FindOneByEmailOrUsername() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(m, nil)

	user, err := s.svc.FindOneByEmailOrUsername(context.Background(), email, username)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), id, user.Id)
	assert.Equal(s.T(), email, user.Email)
	assert.Equal(s.T(), username, user.Username)
}

func (s *UserTestSuite) TestUser_FindOneByEmailOrUsername_Err() {
	email := "test@example.com"
	username := "test"
	m := models.NewUser(email, username)
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOne(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	_, err := s.svc.FindOneByEmailOrUsername(context.Background(), email, username)
	assert.Error(s.T(), err)
	var se *services.Error
	assert.ErrorAs(s.T(), err, &se)
	if errors.As(err, &se) {
		assert.Equal(s.T(), services.NotExist, se.Kind)
	}
}
