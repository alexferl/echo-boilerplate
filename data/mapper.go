package data

import (
	"context"
)

//go:generate mockery --output=../mocks --name Mapper
type Mapper interface {
	Insert(ctx context.Context, document any) error
	FindOne(ctx context.Context, filter any, result any) (any, error)
	FindOneById(ctx context.Context, id string, result any) (any, error)
	Find(ctx context.Context, filter any, result any) (any, error)
	Update(ctx context.Context, filter any, update any, result any) (any, error)
	UpdateById(ctx context.Context, id string, document any, result any) (any, error)
}
