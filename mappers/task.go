package mappers

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
)

// Task represents the mapper used for interacting with Task documents.
type Task struct {
	mapper data.Mapper
}

func NewTask(client *mongo.Client) *Task {
	return &Task{data.NewMapper(client, viper.GetString(config.AppName), "tasks")}
}

func (t *Task) Create(ctx context.Context, model *models.Task) (*models.Task, error) {
	seq, err := t.mapper.GetNextSequence(ctx, "tasks")
	if err != nil {
		return nil, err
	}

	model.Id = seq.String()

	insert, err := t.mapper.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	pipeline := t.getPipeline(bson.D{{"_id", insert.InsertedID.(primitive.ObjectID)}}, 1, 0)
	task, err := t.getTask(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *Task) Find(ctx context.Context, filter any, limit int, skip int) (int64, models.Tasks, error) {
	count, err := t.mapper.Count(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	pipeline := t.getPipeline(filter, limit, skip)
	res, err := t.mapper.Aggregate(ctx, pipeline, models.Tasks{})
	if err != nil {
		return 0, nil, err
	}

	return count, res.(models.Tasks), nil
}

func (t *Task) FindOneById(ctx context.Context, id string) (*models.Task, error) {
	pipeline := t.getPipeline(bson.D{{"id", id}}, 1, 0)
	res, err := t.getTask(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *Task) Update(ctx context.Context, model *models.Task) (*models.Task, error) {
	filter := bson.D{{"id", model.Id}}
	_, err := t.mapper.UpdateOne(ctx, filter, bson.D{{"$set", model}})
	if err != nil {
		return nil, err
	}

	pipeline := t.getPipeline(filter, 1, 0)
	task, err := t.getTask(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *Task) getTask(ctx context.Context, pipeline mongo.Pipeline) (*models.Task, error) {
	res, err := t.mapper.Aggregate(ctx, pipeline, models.Tasks{})
	if err != nil {
		return nil, err
	}

	task := res.(models.Tasks)
	if len(task) < 1 {
		return nil, data.ErrNoDocuments
	}

	return &task[0], nil
}

func (t *Task) getPipeline(filter any, limit int, skip int) mongo.Pipeline {
	if filter == nil {
		filter = bson.D{}
	}

	return mongo.Pipeline{
		{{"$match", filter}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "created_by.id",
			"foreignField": "id",
			"as":           "created_by",
		}}},
		{{"$unwind", "$created_by"}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "updated_by.id",
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
			"localField":   "completed_by.id",
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
