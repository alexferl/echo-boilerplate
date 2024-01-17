package tasks

import (
	"fmt"
	"time"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/users"
	"github.com/alexferl/echo-boilerplate/util"
)

type Task struct {
	*data.Model `bson:",inline"`
	Completed   bool       `json:"completed" bson:"completed"`
	CompletedAt *time.Time `json:"completed_at" bson:"completed_at"`
	CompletedBy string     `json:"completed_by" bson:"completed_by"`
	Title       string     `json:"title" bson:"title"`
}

type Aggregate struct {
	*data.Model `bson:",inline"`
	CreatedAt   *time.Time  `json:"created_at" bson:"created_at"`
	CreatedBy   *users.User `json:"created_by" bson:"created_by"`
	DeletedAt   *time.Time  `json:"-" bson:"deleted_at"`
	DeletedBy   string      `json:"-" bson:"deleted_by"`
	UpdatedAt   *time.Time  `json:"updated_at" bson:"updated_at"`
	UpdatedBy   *users.User `json:"updated_by" bson:"updated_by"`
	Completed   bool        `json:"completed" bson:"completed"`
	CompletedAt *time.Time  `json:"completed_at" bson:"completed_at"`
	CompletedBy *users.User `json:"completed_by" bson:"completed_by"`
	Title       string      `json:"title" bson:"title"`
}

func (a *Aggregate) Response() *Response {
	resp := &Response{
		Id:          a.Id,
		Href:        util.GetFullURL(fmt.Sprintf("/tasks/%s", a.Id)),
		CreatedAt:   a.CreatedAt,
		CreatedBy:   a.CreatedBy.Public(),
		DeletedAt:   a.DeletedAt,
		DeletedBy:   a.DeletedBy,
		UpdatedAt:   a.UpdatedAt,
		Title:       a.Title,
		Completed:   a.Completed,
		CompletedAt: a.CompletedAt,
	}

	if a.UpdatedBy != nil {
		resp.UpdatedBy = a.UpdatedBy.Public()
	}

	if a.CompletedBy != nil {
		resp.CompletedBy = a.CompletedBy.Public()
	}

	return resp
}

type Aggregates []Aggregate

func (aggregates Aggregates) Response() []*Response {
	res := make([]*Response, 0)
	for _, task := range aggregates {
		res = append(res, task.Response())
	}
	return res
}

type Response struct {
	Id          string        `json:"id" bson:"id"`
	Href        string        `json:"href" bson:"href"`
	CreatedAt   *time.Time    `json:"created_at" bson:"created_at"`
	CreatedBy   *users.Public `json:"created_by" bson:"created_by"`
	DeletedAt   *time.Time    `json:"-" bson:"deleted_at"`
	DeletedBy   string        `json:"-" bson:"deleted_by"`
	UpdatedAt   *time.Time    `json:"updated_at" bson:"updated_at"`
	UpdatedBy   *users.Public `json:"updated_by" bson:"updated_by"`
	Completed   bool          `json:"completed" bson:"completed"`
	CompletedAt *time.Time    `json:"completed_at" bson:"completed_at"`
	CompletedBy *users.Public `json:"completed_by" bson:"completed_by"`
	Title       string        `json:"title" bson:"title"`
}

func NewTask(id string) *Task {
	return &Task{
		Model: data.NewModel(id),
	}
}

func (t *Task) Complete(id string) {
	t.Completed = true
	now := time.Now()
	t.CompletedAt = &now
	t.CompletedBy = id
}

func (t *Task) Incomplete() {
	t.Completed = false
	t.CompletedAt = nil
	t.CompletedBy = ""
}

func (t *Task) Aggregate(createdBy *users.User, updatedBy *users.User, completedBy *users.User) *Aggregate {
	resp := &Aggregate{
		Model:       t.Model,
		CreatedAt:   t.CreatedAt,
		CreatedBy:   createdBy,
		DeletedAt:   t.DeletedAt,
		DeletedBy:   t.DeletedBy,
		UpdatedAt:   t.UpdatedAt,
		Title:       t.Title,
		Completed:   t.Completed,
		CompletedAt: t.CompletedAt,
	}

	if updatedBy != nil {
		resp.UpdatedBy = updatedBy
	}

	if completedBy != nil {
		resp.CompletedBy = completedBy
	}

	return resp
}
