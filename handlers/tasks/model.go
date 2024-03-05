package tasks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

type Repository interface {
	New() *Model
	Load(ctx context.Context, id string, token jwt.Token) (*Model, error)
	Find(ctx context.Context, filter any, limit int, skip int) (int64, Tasks, error)

	Create(ctx context.Context, id string, body CreateTaskRequest) (*Task, error)
	Read(ctx context.Context) (*Task, error)
	Update(ctx context.Context, id string) (*Task, error)
	Delete(ctx context.Context, id string) error

	Complete(id string)
	Incomplete()
}

// Model is the base Model that's saved to the database
type Model struct {
	mapper data.Mapper

	*data.Model `bson:",inline"`
	Completed   bool       `json:"completed" bson:"completed"`
	CompletedAt *time.Time `json:"completed_at" bson:"completed_at"`
	CompletedBy *string    `json:"completed_by" bson:"completed_by"`
	Title       string     `json:"title" bson:"title"`
}

type Task struct {
	*Model      `bson:",inline"`
	CreatedBy   *users.User `json:"created_by" bson:"created_by"`
	DeletedAt   *time.Time  `json:"-" bson:"deleted_at"`
	DeletedBy   *string     `json:"-" bson:"deleted_by"`
	UpdatedBy   *users.User `json:"updated_by" bson:"updated_by"`
	CompletedBy *users.User `json:"completed_by" bson:"completed_by"`
}

type Response struct {
	Id          string        `json:"id"`
	Href        string        `json:"href"`
	CreatedAt   *time.Time    `json:"created_at"`
	CreatedBy   *users.Public `json:"created_by"`
	DeletedAt   *time.Time    `json:"-"`
	DeletedBy   *string       `json:"-"`
	UpdatedAt   *time.Time    `json:"updated_at"`
	UpdatedBy   *users.Public `json:"updated_by"`
	Completed   bool          `json:"completed"`
	CompletedAt *time.Time    `json:"completed_at"`
	CompletedBy *users.Public `json:"completed_by"`
	Title       string        `json:"title"`
}

func (t *Task) Response() *Response {
	resp := &Response{
		Id:          t.Id,
		Href:        util.GetFullURL(fmt.Sprintf("/tasks/%s", t.Id)),
		CreatedAt:   t.CreatedAt,
		CreatedBy:   t.CreatedBy.Public(),
		DeletedAt:   t.DeletedAt,
		DeletedBy:   t.DeletedBy,
		UpdatedAt:   t.UpdatedAt,
		Title:       t.Title,
		Completed:   t.Completed,
		CompletedAt: t.CompletedAt,
	}

	if t.UpdatedBy != nil {
		resp.UpdatedBy = t.UpdatedBy.Public()
	}

	if t.CompletedBy != nil {
		resp.CompletedBy = t.CompletedBy.Public()
	}

	return resp
}

type Tasks []Task

func (t Tasks) Response() []*Response {
	res := make([]*Response, 0)
	for _, task := range t {
		res = append(res, task.Response())
	}
	return res
}

func NewModel(mapper data.Mapper) *Model {
	return &Model{Model: data.NewModel(""), mapper: mapper}
}

func (m *Model) New() *Model {
	return m
}

func (m *Model) Load(ctx context.Context, id string, token jwt.Token) (*Model, error) {
	res, err := m.mapper.FindOneById(ctx, id, &Model{})
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, echo.NewHTTPError(http.StatusNotFound, echo.Map{"message": "task not found"})
		}
		return nil, err
	}

	model := res.(*Model)
	model.mapper = m.mapper
	if model.DeletedAt != nil {
		return nil, echo.NewHTTPError(http.StatusGone, echo.Map{"message": "task was deleted"})
	}

	if token != nil {
		if token.Subject() != model.CreatedBy && !util.HasRoles(token, users.AdminRole.String(), users.SuperRole.String()) {
			return nil, echo.NewHTTPError(http.StatusForbidden, echo.Map{"message": "you don't have access"})
		}
	}

	return model, nil
}

func (m *Model) Find(ctx context.Context, filter any, limit int, skip int) (int64, Tasks, error) {
	count, err := m.mapper.Count(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	pipeline := m.getPipeline(filter, limit, skip)
	res, err := m.mapper.Aggregate(ctx, pipeline, Tasks{})
	if err != nil {
		return 0, nil, err
	}

	return count, res.(Tasks), nil
}

func (m *Model) Create(ctx context.Context, id string, body CreateTaskRequest) (*Task, error) {
	m.Model.Create(id)
	m.Title = body.Title

	seq, err := m.mapper.GetNextSequence(ctx, "tasks")
	if err != nil {
		return nil, err
	}

	m.Id = seq.String()

	insert, err := m.mapper.InsertOne(ctx, m)
	if err != nil {
		return nil, err
	}

	pipeline := m.getPipeline(bson.D{{"_id", insert.InsertedID.(primitive.ObjectID)}}, 1, 0)
	task, err := m.getTask(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (m *Model) Read(ctx context.Context) (*Task, error) {
	pipeline := m.getPipeline(bson.D{{"id", m.Id}}, 1, 0)
	task, err := m.getTask(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (m *Model) Update(ctx context.Context, id string) (*Task, error) {
	m.Model.Update(id)

	_, err := m.mapper.UpdateOneById(ctx, m.Id, m)
	if err != nil {
		return nil, err
	}

	pipeline := m.getPipeline(bson.D{{"id", m.Id}}, 1, 0)
	task, err := m.getTask(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (m *Model) Delete(ctx context.Context, id string) error {
	m.Model.Delete(id)

	_, err := m.mapper.UpdateOneById(ctx, id, m)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) Complete(id string) {
	m.Completed = true
	now := time.Now()
	m.CompletedAt = &now
	m.CompletedBy = &id
}

func (m *Model) Incomplete() {
	s := ""
	m.Completed = false
	m.CompletedAt = nil
	m.CompletedBy = &s
}

func (m *Model) getTask(ctx context.Context, pipeline mongo.Pipeline) (*Task, error) {
	res, err := m.mapper.Aggregate(ctx, pipeline, Tasks{})
	if err != nil {
		return nil, err
	}

	resp := res.(Tasks)
	if len(resp) < 1 {
		return nil, errors.New("failed retrieving task")
	}

	return &resp[0], nil
}

func (m *Model) getPipeline(filter any, limit int, skip int) mongo.Pipeline {
	if filter == nil {
		filter = bson.D{}
	}

	return mongo.Pipeline{
		{{"$match", filter}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "created_by",
			"foreignField": "id",
			"as":           "created_by",
		}}},
		{{"$unwind", "$created_by"}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "updated_by",
			"foreignField": "id",
			"as":           "updated_by",
		}}},
		{{
			"$unwind", bson.D{
				{"path", "$updated_by"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "completed_by",
			"foreignField": "id",
			"as":           "completed_by",
		}}},
		{{
			"$unwind", bson.D{
				{"path", "$completed_by"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		{{"$sort", bson.D{{"_id", -1}}}},
		{{"$limit", skip + limit}},
		{{"$skip", skip}},
	}
}
