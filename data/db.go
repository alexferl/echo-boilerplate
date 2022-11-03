package data

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/config"
)

func NewClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := viper.GetString(config.MongoDBURI)
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	opts := options.Client()
	opts.ApplyURI(uri)
	opts.SetAppName(viper.GetString(config.AppName))
	opts.SetServerSelectionTimeout(viper.GetDuration(config.MongoDBServerSelectionTimeoutMs))
	opts.SetConnectTimeout(viper.GetDuration(config.MongoDBConnectTimeoutMs))
	opts.SetSocketTimeout(viper.GetDuration(config.MongoDBSocketTimeoutMs))

	username := viper.GetString(config.MongoDBUsername)
	password := viper.GetString(config.MongoDBPassword)
	if username != "" {
		opts.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})
	}

	replSet := viper.GetString(config.MongoDBReplicaSet)
	if replSet != "" {
		opts.SetReplicaSet(replSet)
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(config.AppName)

	t := true
	_, err = db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{"id", 1},
			},
			Options: &options.IndexOptions{
				Unique: &t,
			},
		},
		{
			Keys: bson.D{
				{"username", 1},
			},
			Options: &options.IndexOptions{
				Unique: &t,
			},
		},
		{
			Keys: bson.D{
				{"email", 1},
			},
			Options: &options.IndexOptions{
				Unique: &t,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	return client, nil
}
