package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexferl/echo-boilerplate/data"
	"github.com/alexferl/echo-boilerplate/models"
)

// UserMapper defines the datastore handling persisting User documents.
type UserMapper interface {
	Create(ctx context.Context, model *models.User) (*models.User, error)
	Find(ctx context.Context, filter any, limit int, skip int) (int64, models.Users, error)
	FindOne(ctx context.Context, filter any) (*models.User, error)
	Update(ctx context.Context, model *models.User) (*models.User, error)
}

var (
	ErrUserDeleted  = errors.New("user was deleted")
	ErrUserExist    = errors.New("email or username already in-use")
	ErrUserNotFound = errors.New("user not found")
)

// User defines the application service in charge of interacting with Users.
type User struct {
	mapper UserMapper
}

func NewUser(mapper UserMapper) *User {
	return &User{mapper: mapper}
}

func (u *User) Create(ctx context.Context, model *models.User) (*models.User, error) {
	model.Create(model.Id)
	res, err := u.mapper.Create(ctx, model)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, NewError(err, Exist, ErrUserExist.Error())
		}
		return nil, NewError(err, Other, "other")
	}

	return res, nil
}

func (u *User) Read(ctx context.Context, id string) (*models.User, error) {
	filter := bson.D{{"$or", bson.A{
		bson.D{{"id", id}},
		bson.D{{"username", id}},
	}}}
	user, err := u.mapper.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, NewError(err, NotExist, ErrUserNotFound.Error())
		}
		return nil, NewError(err, Other, "other")
	}

	if user.DeletedBy != nil {
		return nil, NewError(err, Deleted, ErrUserDeleted.Error())
	}

	return user, nil
}

func (u *User) Update(ctx context.Context, id string, model *models.User) (*models.User, error) {
	// some auth updates aren't 'real' updates
	// and shouldn't update the UpdateAt timestamp
	if id != "" {
		model.Update(id)
	}
	res, err := u.mapper.Update(ctx, model)
	if err != nil {
		return nil, NewError(err, Other, "other")
	}

	return res, nil
}

func (u *User) Delete(ctx context.Context, id string, model *models.User) error {
	model.Delete(id)
	_, err := u.mapper.Update(ctx, model)
	if err != nil {
		return NewError(err, Other, "other")
	}

	return nil
}

func (u *User) Find(ctx context.Context, params *models.UserSearchParams) (int64, models.Users, error) {
	filter := bson.M{"deleted_at": bson.M{"$eq": nil}}
	count, users, err := u.mapper.Find(ctx, filter, params.Limit, params.Skip)
	if err != nil {
		return 0, nil, NewError(err, Other, "other")
	}

	return count, users, nil
}

func (u *User) FindOneByEmailOrUsername(ctx context.Context, email string, username string) (*models.User, error) {
	filter := bson.D{{"$or", bson.A{
		bson.D{{"email", email}},
		bson.D{{"username", username}},
	}}}
	user, err := u.mapper.FindOne(ctx, filter)
	if err != nil {
		if errors.Is(err, data.ErrNoDocuments) {
			return nil, NewError(err, NotExist, ErrUserNotFound.Error())
		}
		return nil, NewError(err, Other, "other")
	}

	return user, nil
}
