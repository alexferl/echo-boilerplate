package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
	"github.com/alexferl/echo-boilerplate/services"
)

type TaskTestSuite struct {
	suite.Suite
	mapper *services.MockTaskMapper
	svc    *services.Task
}

func (s *TaskTestSuite) SetupTest() {
	s.mapper = services.NewMockTaskMapper(s.T())
	s.svc = services.NewTask(s.mapper)
}

func TestTaskTestSuite(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}

func (s *TaskTestSuite) TestTask_Create() {
	m := models.NewTask()
	id := "123"
	m.Create(id)

	s.mapper.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(m, nil)

	task, err := s.svc.Create(context.Background(), id, m)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), task.CreatedBy)
}

func (s *TaskTestSuite) TestTask_Read() {
	m := models.NewTask()
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything).
		Return(m, nil)

	task, err := s.svc.Read(context.Background(), id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), id, task.Id)
}

func (s *TaskTestSuite) TestTask_Read_Err() {
	m := models.NewTask()
	id := "123"
	m.Id = id

	s.mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	_, err := s.svc.Read(context.Background(), id)
	assert.Error(s.T(), err)
	var se *services.Error
	assert.ErrorAs(s.T(), err, &se)
	if errors.As(err, &se) {
		assert.Equal(s.T(), services.NotExist, se.Kind)
	}
}

func (s *TaskTestSuite) TestTask_Update() {
	m := models.NewTask()
	id := "123"
	m.Update(id)

	s.mapper.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(m, nil)

	task, err := s.svc.Update(context.Background(), id, m)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), task.UpdatedBy)
}

func (s *TaskTestSuite) TestTask_Delete() {
	m := models.NewTask()
	id := "123"
	m.Delete(id)

	s.mapper.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(m, nil)

	err := s.svc.Delete(context.Background(), id, m)
	assert.NoError(s.T(), err)

	s.mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything).
		Return(m, nil)

	_, err = s.svc.Read(context.Background(), id)
	assert.Error(s.T(), err)
	var se *services.Error
	assert.ErrorAs(s.T(), err, &se)
	if errors.As(err, &se) {
		assert.Equal(s.T(), services.Deleted, se.Kind)
	}
}

func (s *TaskTestSuite) TestTask_Find() {
	s.mapper.EXPECT().
		Find(mock.Anything, mock.Anything, 1, 0).
		Return(1, models.Tasks{}, nil)

	count, tasks, err := s.svc.Find(context.Background(), &models.TaskSearchParams{
		Completed: []string{"true", "false"},
		CreatedBy: "123",
		Queries:   []string{"foo", "bar"},
		Limit:     1,
		Skip:      0,
	})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), count)
	assert.Equal(s.T(), models.Tasks{}, tasks)
}
