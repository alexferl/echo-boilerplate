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

// PersonalAccessToken represents the mapper used for interacting with PersonalAccessToken documents.
type PersonalAccessToken struct {
	mapper data.Mapper
}

func NewPersonalAccessToken(client *mongo.Client) *PersonalAccessToken {
	return &PersonalAccessToken{data.NewMapper(client, viper.GetString(config.AppName), "personal_access_tokens")}
}

func (p PersonalAccessToken) Create(ctx context.Context, model *models.PersonalAccessToken) (*models.PersonalAccessToken, error) {
	filter := bson.D{{"id", model.Id}}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	res, err := p.mapper.FindOneAndUpdate(ctx, filter, model, &models.PersonalAccessToken{}, opts)
	if err != nil {
		return nil, err
	}

	return res.(*models.PersonalAccessToken), nil
}

func (p PersonalAccessToken) Find(ctx context.Context, filter any) (models.PersonalAccessTokens, error) {
	res, err := p.mapper.Find(ctx, filter, models.PersonalAccessTokens{})
	if err != nil {
		return nil, err
	}

	return res.(models.PersonalAccessTokens), nil
}

func (p PersonalAccessToken) FindOne(ctx context.Context, filter any) (*models.PersonalAccessToken, error) {
	res, err := p.mapper.FindOne(ctx, filter, &models.PersonalAccessToken{})
	if err != nil {
		return nil, err
	}

	return res.(*models.PersonalAccessToken), nil
}

func (p PersonalAccessToken) Update(ctx context.Context, model *models.PersonalAccessToken) (*models.PersonalAccessToken, error) {
	filter := bson.D{{"id", model.Id}}
	res, err := p.mapper.FindOneAndUpdate(ctx, filter, model, &models.PersonalAccessToken{})
	if err != nil {
		return nil, err
	}

	return res.(*models.PersonalAccessToken), nil
}
