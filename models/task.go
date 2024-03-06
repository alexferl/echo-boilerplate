package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alexferl/echo-boilerplate/util"
)

type Task struct {
	Model       `bson:",inline"`
	Completed   bool       `json:"completed" bson:"completed"`
	CompletedAt *time.Time `json:"completed_at" bson:"completed_at"`
	CompletedBy any        `json:"completed_by" bson:"completed_by"`
	Title       string     `json:"title" bson:"title"`
}

type Tasks []Task

func (t Tasks) Response() []TaskResponse {
	res := make([]TaskResponse, 0)
	for _, task := range t {
		res = append(res, task.Response())
	}
	return res
}

func NewTask() *Task {
	return &Task{Model: NewModel()}
}

type TaskResponse struct {
	Id          string     `json:"id"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *Public    `json:"created_by"`
	DeletedAt   *time.Time `json:"-"`
	DeletedBy   string     `json:"-"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *Public    `json:"updated_by"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	CompletedBy *Public    `json:"completed_by"`
	Title       string     `json:"title"`
}

func (t *Task) Response() TaskResponse {
	resp := TaskResponse{
		Id:          t.Id,
		CreatedAt:   t.CreatedAt,
		CreatedBy:   t.CreatedBy.(*User).Public(),
		UpdatedAt:   t.UpdatedAt,
		Title:       t.Title,
		Completed:   t.Completed,
		CompletedAt: t.CompletedAt,
	}

	if t.CompletedBy != nil {
		resp.CompletedBy = t.CompletedBy.(*User).Public()
	}

	if t.UpdatedBy != nil {
		resp.UpdatedBy = t.UpdatedBy.(*User).Public()
	}

	return resp
}

func (t *Task) Complete(id string) {
	t.Completed = true
	now := time.Now()
	t.CompletedAt = &now
	t.CompletedBy = &Ref{Id: id}
}

func (t *Task) Incomplete() {
	t.Completed = false
	t.CompletedAt = nil
	t.CompletedBy = nil
}

func (t *Task) UnmarshalBSON(b []byte) error {
	type Alias Task

	if err := bson.Unmarshal(b, (*Alias)(t)); err != nil {
		return err
	}

	if t.CompletedBy != nil {
		var u *User
		err := util.DocToStruct(t.CompletedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.CompletedBy = u
	}

	if t.CreatedBy != nil {
		var u *User
		err := util.DocToStruct(t.CreatedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.CreatedBy = u
	}

	if t.DeletedBy != nil {
		var u *User
		err := util.DocToStruct(t.DeletedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.DeletedBy = u
	}

	if t.UpdatedBy != nil {
		var u *User
		err := util.DocToStruct(t.UpdatedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.UpdatedBy = u
	}

	return nil
}

type TaskSearchParams struct {
	Completed []string
	CreatedBy string
	Queries   []string
	Limit     int
	Skip      int
}
