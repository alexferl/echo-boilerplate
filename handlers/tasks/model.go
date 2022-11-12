package tasks

import (
	"time"

	"github.com/alexferl/echo-boilerplate/data"
)

type Task struct {
	*data.Model `bson:",inline"`
	Title       string     `json:"title" bson:"title"`
	IsPrivate   bool       `json:"is_private" bson:"is_private"`
	IsCompleted bool       `json:"is_completed" bson:"is_completed"`
	CompletedAt *time.Time `json:"completed_at" bson:"completed_at"`
	CompletedBy string     `json:"completed_by" bson:"completed_by"`
}

type TaskWithUsers struct {
	Id          string     `json:"id" bson:"id"`
	CreatedAt   *time.Time `json:"created_at" bson:"created_at"`
	CreatedBy   *TaskUser  `json:"created_by" bson:"created_by"`
	DeletedAt   *time.Time `json:"deleted_at" bson:"deleted_at"`
	DeletedBy   string     `json:"deleted_by" bson:"deleted_by"`
	UpdatedAt   *time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy   *TaskUser  `json:"updated_by" bson:"updated_by"`
	Title       string     `json:"title" bson:"title"`
	IsPrivate   bool       `json:"is_private" bson:"is_private"`
	IsCompleted bool       `json:"is_completed" bson:"is_completed"`
	CompletedAt *time.Time `json:"completed_at" bson:"completed_at"`
	CompletedBy *TaskUser  `json:"completed_by" bson:"completed_by"`
}

type TaskUser struct {
	Id       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
}

func NewTask() *Task {
	return &Task{
		Model: data.NewModel(),
	}
}

func (t *Task) SetPrivate() {
	t.IsPrivate = true
}

func (t *Task) Complete(id string) {
	t.IsCompleted = true
	now := time.Now()
	t.CompletedAt = &now
	t.CompletedBy = id
}
