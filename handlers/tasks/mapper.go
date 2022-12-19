package tasks

import (
	"context"
	"errors"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/alexferl/echo-boilerplate/config"
	"github.com/alexferl/echo-boilerplate/data"
)

var ErrTaskNotFound = errors.New("task not found")

type Mapper struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMapper(db *mongo.Client) data.Mapper {
	collection := db.Database(viper.GetString(config.AppName)).Collection("tasks")
	return &Mapper{
		db,
		collection,
	}
}

func (m *Mapper) Collection(name string) data.Mapper {
	// TODO implement me
	panic("implement me")
}

func (m *Mapper) Insert(ctx context.Context, document any, result any, opts ...*options.InsertOneOptions) (any, error) {
	session, txnOpts, err := getSession(m.client)
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongo.SessionContext) (any, error) {
		res, err := m.collection.InsertOne(sessionCtx, document, opts...)
		if err != nil {
			return nil, err
		}

		filter := bson.D{{"_id", res.InsertedID}}
		return m.Aggregate(sessionCtx, filter, 1, 0, result)
	}

	res, err := session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Mapper) FindOne(ctx context.Context, filter any, result any, opts ...*options.FindOneOptions) (any, error) {
	err := m.collection.FindOne(ctx, filter, opts...).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil, ErrTaskNotFound
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Mapper) FindOneById(ctx context.Context, id string, result any, opts ...*options.FindOneOptions) (any, error) {
	return m.FindOne(ctx, bson.D{{"id", id}}, result, opts...)
}

func (m *Mapper) Find(ctx context.Context, filter any, result any, opts ...*options.FindOptions) (any, error) {
	if filter == nil {
		filter = bson.D{}
	}

	cur, err := m.collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Mapper) Aggregate(ctx context.Context, filter any, limit int, skip int, result any, opts ...*options.AggregateOptions) (any, error) {
	if filter == nil {
		filter = bson.D{}
	}
	pipeline := mongo.Pipeline{
		{{"$match", filter}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "created_by",
			"foreignField": "id",
			"as":           "created_by",
		}}},
		{{"$unwind", "$created_by"}},
		{{"$lookup", bson.M{
			"from":         "users",
			"localField":   "updated_by",
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
			"localField":   "completed_by",
			"foreignField": "id",
			"as":           "completed_by",
		}}},
		{{
			"$unwind", bson.D{
				{"path", "$completed_by"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		{{"$limit", skip + limit}},
		{{"$skip", skip}},
	}

	cur, err := m.collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Mapper) Count(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	if filter == nil {
		filter = bson.D{}
	}

	count, err := m.collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *Mapper) Update(ctx context.Context, filter any, update any, result any, opts ...*options.UpdateOptions) (any, error) {
	session, txnOpts, err := getSession(m.client)
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	callback := func(sessionCtx mongo.SessionContext) (any, error) {
		_, err = m.collection.UpdateOne(sessionCtx, filter, update, opts...)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return m.Aggregate(sessionCtx, filter, 1, 0, result)
		}

		return nil, nil
	}

	res, err := session.WithTransaction(context.Background(), callback, txnOpts)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Mapper) UpdateById(ctx context.Context, id string, document any, result any, opts ...*options.UpdateOptions) (any, error) {
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", document}}

	return m.Update(ctx, filter, update, result, opts...)
}

func (m *Mapper) Upsert(ctx context.Context, filter any, update any, result any, opts ...*options.FindOneAndUpdateOptions) (any, error) {
	// TODO implement me
	panic("implement me")
}

func getSession(client *mongo.Client) (mongo.Session, *options.TransactionOptions, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	session, err := client.StartSession()
	if err != nil {
		return nil, nil, err
	}

	return session, txnOpts, nil
}
