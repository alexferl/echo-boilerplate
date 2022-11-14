package data

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockery --output=../mocks --name Mapper
type Mapper interface {
	Collection(name string) Mapper
	Insert(ctx context.Context, document any, result any, opts ...*options.InsertOneOptions) (any, error)
	FindOne(ctx context.Context, filter any, result any, opts ...*options.FindOneOptions) (any, error)
	FindOneById(ctx context.Context, id string, result any, opts ...*options.FindOneOptions) (any, error)
	Find(ctx context.Context, filter any, result any, opts ...*options.FindOptions) (any, error)
	Aggregate(ctx context.Context, filter any, limit int, skip int, result any, opts ...*options.AggregateOptions) (any, error)
	Count(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)
	Update(ctx context.Context, filter any, update any, result any, opts ...*options.UpdateOptions) (any, error)
	UpdateById(ctx context.Context, id string, document any, result any, opts ...*options.UpdateOptions) (any, error)
	Upsert(ctx context.Context, filter any, update any, result any, opts ...*options.FindOneAndUpdateOptions) (any, error)
}
