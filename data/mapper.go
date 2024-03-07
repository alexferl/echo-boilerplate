package data

import (
	"context"
	"errors"
	"strconv"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

var ErrNoDocuments = errors.New("no documents in result")

type Mapper interface {
	Aggregate(ctx context.Context, pipeline mongo.Pipeline, results any, opts ...*options.AggregateOptions) (any, error)
	Count(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)
	Find(ctx context.Context, filter any, results any, opts ...*options.FindOptions) (any, error)
	FindOne(ctx context.Context, filter any, result any, opts ...*options.FindOneOptions) (any, error)
	FindOneAndUpdate(ctx context.Context, filter any, update any, result any, opts ...*options.FindOneAndUpdateOptions) (any, error)
	InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	GetNextSequence(ctx context.Context, name string) (*Sequence, error)
}

type mapper struct {
	client     *mongo.Client
	db         *mongo.Database
	dbName     string
	collection *mongo.Collection
}

func NewMapper(client *mongo.Client, databaseName string, collectionName string) Mapper {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	return &mapper{
		client,
		db,
		databaseName,
		collection,
	}
}

func (m *mapper) Aggregate(ctx context.Context, pipeline mongo.Pipeline, results any, opts ...*options.AggregateOptions) (any, error) {
	cur, err := m.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msg("cursor error")
		}
	}(cur, ctx)

	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mapper) Count(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := m.collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *mapper) Find(ctx context.Context, filter any, results any, opts ...*options.FindOptions) (any, error) {
	if filter == nil {
		filter = bson.D{}
	}

	cur, err := m.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error().Err(err).Msg("cursor error")
		}
	}(cur, ctx)

	err = cur.All(ctx, &results)
	if err != nil {
		return nil, err
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mapper) FindOne(ctx context.Context, filter any, result any, opts ...*options.FindOneOptions) (any, error) {
	err := m.collection.FindOne(ctx, filter, opts...).Decode(result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNoDocuments
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *mapper) FindOneAndUpdate(ctx context.Context, filter any, update any, result any, opts ...*options.FindOneAndUpdateOptions) (any, error) {
	opts = append(opts, options.FindOneAndUpdate().SetReturnDocument(options.After))
	res := m.collection.FindOneAndUpdate(ctx, filter, bson.D{{"$set", update}}, opts...)
	if res.Err() != nil {
		return nil, res.Err()
	}

	if result != nil {
		err := res.Decode(result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (m *mapper) InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	res, err := m.collection.InsertOne(ctx, document, opts...)
	return res, err
}

func (m *mapper) UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	res, err := m.collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Sequence struct {
	Seq int `bson:"seq"`
}

func (s *Sequence) String() string {
	return strconv.Itoa(s.Seq)
}

func (m *mapper) GetNextSequence(ctx context.Context, name string) (*Sequence, error) {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	res := m.db.Collection("counters").FindOneAndUpdate(
		ctx,
		bson.D{{"_id", name}},
		bson.D{{"$inc", bson.D{{"seq", 1}}}},
		opts,
	)
	if res.Err() != nil {
		return nil, res.Err()
	}

	seq := &Sequence{}
	err := res.Decode(seq)
	if err != nil {
		return nil, err
	}

	return seq, nil
}

func (m *mapper) getSession() (mongo.Session, *options.TransactionOptions, error) {
	wc := writeconcern.Majority()
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := m.client.StartSession()
	if err != nil {
		return nil, nil, err
	}

	return session, txnOpts, nil
}
