package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
)

// PersonalAccessTokenMapper defines the datastore handling persisting User documents.
type PersonalAccessTokenMapper interface {
	Create(ctx context.Context, model *models.PersonalAccessToken) (*models.PersonalAccessToken, error)
	Find(ctx context.Context, filter any) (models.PersonalAccessTokens, error)
	FindOne(ctx context.Context, filter any) (*models.PersonalAccessToken, error)
	Update(ctx context.Context, model *models.PersonalAccessToken) (*models.PersonalAccessToken, error)
}

var ErrPersonalAccessTokenNotFound = errors.New("personal access token not found")

// PersonalAccessToken defines the application service in charge of interacting with Users.
type PersonalAccessToken struct {
	mapper PersonalAccessTokenMapper
}

func NewPersonalAccessToken(mapper PersonalAccessTokenMapper) *PersonalAccessToken {
	return &PersonalAccessToken{mapper: mapper}
}

func (t *PersonalAccessToken) Create(ctx context.Context, model *models.PersonalAccessToken) (*models.PersonalAccessToken, error) {
	token, err := t.mapper.Create(ctx, model)
	if err != nil {
		return nil, NewError(err, Other, "other")
	}

	return token, nil
}

func (t *PersonalAccessToken) Read(ctx context.Context, userId string, id string) (*models.PersonalAccessToken, error) {
	filter := bson.D{{"user_id", userId}, {"id", id}}
	token, err := t.mapper.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, NewError(err, NotExist, ErrPersonalAccessTokenNotFound.Error())
		}
		return nil, NewError(err, Other, "other")
	}

	return token, nil
}

func (t *PersonalAccessToken) Revoke(ctx context.Context, model *models.PersonalAccessToken) error {
	model.IsRevoked = true
	_, err := t.mapper.Update(ctx, model)
	if err != nil {
		return NewError(err, Other, "other")
	}

	return nil
}

func (t *PersonalAccessToken) Find(ctx context.Context, userId string) (models.PersonalAccessTokens, error) {
	filter := bson.D{{"user_id", userId}}
	tokens, err := t.mapper.Find(ctx, filter)
	if err != nil {
		return nil, NewError(err, Other, "other")
	}

	return tokens, nil
}

func (t *PersonalAccessToken) FindOne(ctx context.Context, userId string, name string) (*models.PersonalAccessToken, error) {
	filter := bson.D{{"user_id", userId}, {"name", name}}
	user, err := t.mapper.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, NewError(err, NotExist, ErrPersonalAccessTokenNotFound.Error())
		}
		return nil, NewError(err, Other, "other")
	}

	return user, nil
}
