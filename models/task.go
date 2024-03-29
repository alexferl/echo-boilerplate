package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	utilBSON "github.com/alexferl/echo-boilerplate/util/bson"
)

type Task struct {
	*Model      `bson:",inline"`
	Completed   bool       `bson:"completed"`
	CompletedAt *time.Time `bson:"completed_at"`
	CompletedBy any        `bson:"completed_by"`
	Title       string     `bson:"title"`
}

type TaskResponse struct {
	Id          string     `json:"id"`
	Completed   bool       `json:"completed"`
	CompletedAt *time.Time `json:"completed_at"`
	CompletedBy *UserRef   `json:"completed_by"`
	CreatedAt   *time.Time `json:"created_at"`
	CreatedBy   *UserRef   `json:"created_by"`
	DeletedAt   *time.Time `json:"-"`
	DeletedBy   *UserRef   `json:"-"`
	Title       string     `json:"title"`
	UpdatedAt   *time.Time `json:"updated_at"`
	UpdatedBy   *UserRef   `json:"updated_by"`
}

func NewTask() *Task {
	return &Task{Model: NewModel()}
}

func (t *Task) Response() *TaskResponse {
	resp := &TaskResponse{
		Id:          t.Id,
		Completed:   t.Completed,
		CompletedAt: t.CompletedAt,
		CreatedAt:   t.CreatedAt,
		CreatedBy:   t.CreatedBy.(*User).Ref(),
		Title:       t.Title,
		UpdatedAt:   t.UpdatedAt,
	}

	if t.CompletedBy != nil {
		resp.CompletedBy = t.CompletedBy.(*User).Ref()
	}

	if t.UpdatedBy != nil {
		resp.UpdatedBy = t.UpdatedBy.(*User).Ref()
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

func (t *Task) MarshalBSON() ([]byte, error) {
	type Alias Task
	aux := &struct {
		*Alias `bson:",inline"`
	}{
		Alias: (*Alias)(t),
	}

	if t.CompletedBy != nil {
		user, ok := t.CompletedBy.(*User)
		if ok {
			aux.CompletedBy = &Ref{Id: user.Id}
		}
	}

	if t.CreatedBy != nil {
		user, ok := t.CreatedBy.(*User)
		if ok {
			aux.CreatedBy = &Ref{Id: user.Id}
		}
	}

	if t.DeletedBy != nil {
		user, ok := t.DeletedBy.(*User)
		if ok {
			aux.DeletedBy = &Ref{Id: user.Id}
		}
	}

	if t.UpdatedBy != nil {
		user, ok := t.UpdatedBy.(*User)
		if ok {
			aux.UpdatedBy = &Ref{Id: user.Id}
		}
	}

	return bson.Marshal(aux)
}

func (t *Task) UnmarshalBSON(data []byte) error {
	type Alias Task
	aux := &struct {
		*Alias `bson:",inline"`
	}{
		Alias: (*Alias)(t),
	}

	if err := bson.Unmarshal(data, aux); err != nil {
		return err
	}

	if t.CompletedBy != nil {
		var u *User
		err := utilBSON.DocToStruct(aux.CompletedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.CompletedBy = u
	}

	if t.CreatedBy != nil {
		var u *User
		err := utilBSON.DocToStruct(aux.CreatedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.CreatedBy = u
	}

	if t.DeletedBy != nil {
		var u *User
		err := utilBSON.DocToStruct(aux.DeletedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.DeletedBy = u
	}

	if t.UpdatedBy != nil {
		var u *User
		err := utilBSON.DocToStruct(aux.UpdatedBy.(primitive.D), &u)
		if err != nil {
			return err
		}
		t.UpdatedBy = u
	}

	return nil
}

type Tasks []Task

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

func (t Tasks) Response() *TasksResponse {
	res := make([]TaskResponse, 0)
	for _, task := range t {
		res = append(res, *task.Response())
	}
	return &TasksResponse{Tasks: res}
}

type TaskSearchParams struct {
	Completed []string
	CreatedBy string
	Queries   []string
	Limit     int
	Skip      int
}
