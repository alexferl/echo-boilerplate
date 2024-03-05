package tasks

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/data"
	_ "github.com/alexferl/echo-boilerplate/testing"
)

func Test_Model_Load(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{}}, nil)

	_, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)
}

func Test_Model_Load_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	_, err := m.Load(context.Background(), "", nil)
	assert.Error(t, err)
}

func Test_Model_Load_Not_Found(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, data.ErrNoDocuments)

	_, err := m.Load(context.Background(), "", nil)
	assert.Error(t, err)
	var he *echo.HTTPError
	if errors.As(err, &he) {
		assert.Equal(t, http.StatusNotFound, he.Code)
	}
}

func Test_Model_Load_Not_Gone(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	now := time.Now()

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{DeletedAt: &now}}, nil)

	_, err := m.Load(context.Background(), "", nil)
	assert.Error(t, err)
	var he *echo.HTTPError
	if errors.As(err, &he) {
		assert.Equal(t, http.StatusGone, he.Code)
	}
}

func Test_Model_Load_Not_Forbidden(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	claims := map[string]any{
		"roles": []string{"user"},
	}
	builder := jwt.NewBuilder().Subject("1")
	for k, v := range claims {
		builder.Claim(k, v)
	}
	token, err := builder.Build()
	assert.NoError(t, err)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{}}, nil)

	_, err = m.Load(context.Background(), "", token)
	assert.Error(t, err)
	var he *echo.HTTPError
	if errors.As(err, &he) {
		assert.Equal(t, http.StatusForbidden, he.Code)
	}
}

func Test_Model_Find(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{}}, nil)

	model, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		Count(mock.Anything, mock.Anything).
		Return(1, nil)
	mapper.EXPECT().
		Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(Tasks{{}}, nil)

	num, tasks, err := model.Find(context.Background(), bson.D{}, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), num)
	assert.Len(t, tasks, 1)
}

func Test_Model_Find_Count_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{}}, nil)

	model, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		Count(mock.Anything, mock.Anything).
		Return(0, errors.New(""))

	_, _, err = model.Find(context.Background(), bson.D{}, 1, 0)
	assert.Error(t, err)
}

func Test_Model_Find_Aggregate_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{}}, nil)

	model, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		Count(mock.Anything, mock.Anything).
		Return(1, nil)
	mapper.EXPECT().
		Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	_, _, err = model.Find(context.Background(), bson.D{}, 1, 0)
	assert.Error(t, err)
}

func Test_Model_Create(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	model := m.New()

	mapper.EXPECT().
		GetNextSequence(mock.Anything, mock.Anything).
		Return(&data.Sequence{Seq: 1}, nil)
	mapper.EXPECT().
		InsertOne(mock.Anything, model).
		Return(&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil)
	mapper.EXPECT().Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(Tasks{{}}, nil)

	_, err := model.Create(context.Background(), "", CreateTaskRequest{})
	assert.NoError(t, err)
}

func Test_Model_Create_GetNextSequence_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	model := m.New()

	mapper.EXPECT().
		GetNextSequence(mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	_, err := model.Create(context.Background(), "", CreateTaskRequest{})
	assert.Error(t, err)
}

func Test_Model_Create_InsertOne_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	model := m.New()

	mapper.EXPECT().
		GetNextSequence(mock.Anything, mock.Anything).
		Return(&data.Sequence{Seq: 1}, nil)
	mapper.EXPECT().
		InsertOne(mock.Anything, model).
		Return(nil, errors.New(""))

	_, err := model.Create(context.Background(), "", CreateTaskRequest{})
	assert.Error(t, err)
}

func Test_Model_Create_Aggregate_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	model := m.New()

	mapper.EXPECT().
		GetNextSequence(mock.Anything, mock.Anything).
		Return(&data.Sequence{Seq: 1}, nil)
	mapper.EXPECT().
		InsertOne(mock.Anything, model).
		Return(&mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil)
	mapper.EXPECT().Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	_, err := model.Create(context.Background(), "", CreateTaskRequest{})
	assert.Error(t, err)
}

func Test_Model_Read(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := "1"

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{Id: id}}, nil)

	model, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(Tasks{{Model: &Model{Model: &data.Model{Id: id}}}}, nil)

	task, err := model.Read(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, id, task.Id)
}

func Test_Model_Read_Aggregate_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{}}, nil)

	model, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	_, err = model.Read(context.Background())
	assert.Error(t, err)
}

func Test_Model_Update(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := "1"

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{Id: id}}, nil)

	model, err := m.Load(context.Background(), id, nil)
	assert.NoError(t, err)

	title := "test"

	mapper.EXPECT().
		UpdateOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
	mapper.EXPECT().
		Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(Tasks{{Model: &Model{Model: &data.Model{Id: id}, Title: title}}}, nil)

	task, err := model.Update(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, task.Id)
	assert.Equal(t, title, task.Title)
}

func Test_Model_Update_UpdateOneById_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := "1"

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{Id: id}}, nil)

	model, err := m.Load(context.Background(), id, nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		UpdateOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	_, err = model.Update(context.Background(), id)
	assert.Error(t, err)
}

func Test_Model_Update_Aggregate_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := "1"

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{Id: id}}, nil)

	model, err := m.Load(context.Background(), id, nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		UpdateOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
	mapper.EXPECT().
		Aggregate(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	_, err = model.Update(context.Background(), id)
	assert.Error(t, err)
}

func Test_Model_Delete(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := "1"

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{Id: id}}, nil)

	model, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		UpdateOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)

	err = model.Delete(context.Background(), id)
	assert.NoError(t, err)
}

func Test_Model_Delete_UpdateOneById_Err(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := "1"

	mapper.EXPECT().
		FindOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(&Model{Model: &data.Model{Id: id}}, nil)

	model, err := m.Load(context.Background(), "", nil)
	assert.NoError(t, err)

	mapper.EXPECT().
		UpdateOneById(mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New(""))

	err = model.Delete(context.Background(), id)
	assert.Error(t, err)
}

func Test_Model_Complete(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := "1"
	m.Complete(id)

	assert.Equal(t, &id, m.CompletedBy)
}

func Test_Model_Incomplete(t *testing.T) {
	mapper := data.NewMockMapper(t)
	m := NewModel(mapper)
	id := ""
	m.Incomplete()

	assert.Equal(t, &id, m.CompletedBy)
}
