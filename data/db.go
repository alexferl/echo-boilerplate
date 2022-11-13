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

	db := client.Database(viper.GetString(config.AppName))

	idxOpts := options.Index().
		SetUnique(true).
		SetCollation(&options.Collation{Locale: "en", Strength: 2})

	usernameOpts := idxOpts.SetName("username")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"username", 1}},
		Options: usernameOpts,
	}
	_, err = db.Collection("users").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic(err)
	}

	emailOpts := idxOpts.SetName("email")
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: emailOpts,
	}
	_, err = db.Collection("users").Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		panic(err)
	}

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
	})
	if err != nil {
		panic(err)
	}

	_, err = db.Collection("tasks").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{"id", 1},
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
