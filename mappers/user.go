package mappers

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
)

type User struct {
	mapper data.Mapper
}

func NewUser(client *mongo.Client) *User {
	return &User{data.NewMapper(client, viper.GetString(config.AppName), "users")}
}

func (u *User) Create(ctx context.Context, model *models.User) (*models.User, error) {
	filter := bson.D{{"id", model.Id}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	res, err := u.mapper.FindOneAndUpdate(ctx, filter, model, &models.User{}, opts)
	if err != nil {
		return nil, err
	}

	return res.(*models.User), nil
}

func (u *User) Find(ctx context.Context, filter any, limit int, skip int) (int64, models.Users, error) {
	count, err := u.mapper.Count(ctx, filter)
	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	res, err := u.mapper.Find(ctx, filter, models.Users{}, opts)
	if err != nil {
		return 0, nil, err
	}

	return count, res.(models.Users), nil
}

func (u *User) FindOne(ctx context.Context, filter any) (*models.User, error) {
	res, err := u.mapper.FindOne(ctx, filter, &models.User{})
	if err != nil {
		return nil, err
	}

	return res.(*models.User), nil
}

func (u *User) Update(ctx context.Context, model *models.User) (*models.User, error) {
	filter := bson.D{{"id", model.Id}}
	res, err := u.mapper.FindOneAndUpdate(ctx, filter, model, &models.User{})
	if err != nil {
		return nil, err
	}

	return res.(*models.User), nil
}
