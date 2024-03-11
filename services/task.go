package services

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
)

// TaskMapper defines the datastore handling persisting Task documents.
type TaskMapper interface {
	Create(ctx context.Context, model *models.Task) (*models.Task, error)
	Find(ctx context.Context, filter any, limit int, skip int) (int64, models.Tasks, error)
	FindOneById(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, model *models.Task) (*models.Task, error)
}

var (
	ErrTaskDeleted  = errors.New("task was deleted")
	ErrTaskNotFound = errors.New("task not found")
)

// Task defines the application service in charge of interacting with Tasks.
type Task struct {
	mapper TaskMapper
}

func NewTask(mapper TaskMapper) *Task {
	return &Task{mapper: mapper}
}

func (t *Task) Create(ctx context.Context, id string, model *models.Task) (*models.Task, error) {
	model.Create(id)
	task, err := t.mapper.Create(ctx, model)
	if err != nil {
		return nil, NewError(err, Other, "other")
	}

	return task, nil
}

func (t *Task) Read(ctx context.Context, id string) (*models.Task, error) {
	task, err := t.mapper.FindOneById(ctx, id)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, NewError(err, NotExist, ErrTaskNotFound.Error())
		}
		return nil, NewError(err, Other, "other")
	}

	if task.DeletedBy != nil {
		return nil, NewError(err, Deleted, ErrTaskDeleted.Error())
	}

	return task, nil
}

func (t *Task) Update(ctx context.Context, id string, model *models.Task) (*models.Task, error) {
	model.Update(id)
	task, err := t.mapper.Update(ctx, model)
	if err != nil {
		return nil, NewError(err, Other, "other")
	}

	return task, nil
}

func (t *Task) Delete(ctx context.Context, id string, model *models.Task) error {
	model.Delete(id)
	_, err := t.mapper.Update(ctx, model)
	if err != nil {
		return NewError(err, Other, "other")
	}

	return nil
}

func (t *Task) Find(ctx context.Context, params *models.TaskSearchParams) (int64, models.Tasks, error) {
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	completed := params.Completed
	if len(completed) > 0 {
		arr := bson.A{}
		for _, i := range completed {
			s := strings.ToLower(i)
			if s == "true" {
				arr = append(arr, true)
			} else if s == "false" {
				arr = append(arr, false)
			}
		}
		filter["completed"] = bson.M{"$in": arr}
	}
	createdBy := params.CreatedBy
	if createdBy != "" {
		filter["created_by"] = createdBy
	}
	query := params.Queries
	if len(query) > 0 {
		filter["$text"] = bson.M{"$search": strings.Join(query, " ")}
	}

	count, tasks, err := t.mapper.Find(ctx, filter, params.Limit, params.Skip)
	if err != nil {
		return 0, nil, NewError(err, Other, "other")
	}

	return count, tasks, nil
}
