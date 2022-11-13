package tasks

import (
	"time"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/handlers/users"
)

type Task struct {
	*data.Model `bson:",inline"`
	Title       string     `json:"title" bson:"title"`
	IsCompleted bool       `json:"is_completed" bson:"is_completed"`
	CompletedAt *time.Time `json:"completed_at" bson:"completed_at"`
	CompletedBy string     `json:"completed_by" bson:"completed_by"`
}

func NewTask() *Task {
	return &Task{
		Model: data.NewModel(),
	}
}

func (t *Task) Complete(id string) {
	t.IsCompleted = true
	now := time.Now()
	t.CompletedAt = &now
	t.CompletedBy = id
}

func (t *Task) Incomplete() {
	t.IsCompleted = false
	t.CompletedAt = nil
	t.CompletedBy = ""
}

type TaskResponse struct {
	Id          string            `json:"id" bson:"id"`
	CreatedAt   *time.Time        `json:"created_at" bson:"created_at"`
	CreatedBy   *users.PublicUser `json:"created_by" bson:"created_by"`
	DeletedAt   *time.Time        `json:"-" bson:"deleted_at"`
	DeletedBy   string            `json:"-" bson:"deleted_by"`
	UpdatedAt   *time.Time        `json:"updated_at" bson:"updated_at"`
	UpdatedBy   *users.PublicUser `json:"updated_by" bson:"updated_by"`
	Title       string            `json:"title" bson:"title"`
	IsCompleted bool              `json:"is_completed" bson:"is_completed"`
	CompletedAt *time.Time        `json:"completed_at" bson:"completed_at"`
	CompletedBy *users.PublicUser `json:"completed_by" bson:"completed_by"`
}

func (t *Task) MakeResponse(createdBy *users.User, updatedBy *users.User, completedBy *users.User) *TaskResponse {
	resp := &TaskResponse{
		Id:        t.Id,
		CreatedAt: t.CreatedAt,
		CreatedBy: &users.PublicUser{
			Id:       createdBy.Id,
			Username: createdBy.Username,
			Name:     createdBy.Name,
		},
		DeletedAt:   t.DeletedAt,
		DeletedBy:   t.DeletedBy,
		UpdatedAt:   t.UpdatedAt,
		Title:       t.Title,
		IsCompleted: t.IsCompleted,
		CompletedAt: t.CompletedAt,
	}

	if updatedBy != nil {
		resp.UpdatedBy = &users.PublicUser{
			Id:       updatedBy.Id,
			Username: updatedBy.Username,
			Name:     updatedBy.Name,
		}
	}

	if completedBy != nil {
		resp.CompletedBy = &users.PublicUser{
			Id:       completedBy.Id,
			Username: completedBy.Username,
			Name:     completedBy.Name,
		}
	}

	return resp
}
