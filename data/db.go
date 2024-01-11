package data

import (
	"context"
	"time"

	"github.com/alexferl/golib/database/mongodb"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/alexferl/echo-boilerplate/config"
)

func CreateIndexes(client *mongo.Client) error {
	indexes := map[string][]mongo.IndexModel{}

	username := "username"
	email := "email"
	t := true

	indexes["users"] = []mongo.IndexModel{
		{
			Keys: bson.D{{"username", 1}},
			Options: &options.IndexOptions{
				Name:      &username,
				Unique:    &t,
				Collation: &options.Collation{Locale: "en", Strength: 2},
			},
		},
		{
			Keys: bson.D{{"email", 1}},
			Options: &options.IndexOptions{
				Name:      &email,
				Unique:    &t,
				Collation: &options.Collation{Locale: "en", Strength: 2},
			},
		},
		{
			Keys: bson.D{
				{"id", 1},
			},
			Options: &options.IndexOptions{
				Unique: &t,
			},
		},
	}

	indexes["tasks"] = []mongo.IndexModel{
		{
			Keys: bson.D{
				{"id", 1},
			},
			Options: &options.IndexOptions{
				Unique: &t,
			},
		},
	}

	indexes["personal_access_tokens"] = []mongo.IndexModel{
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
				{"user_id", 1},
			},
		},
		{
			Keys: bson.D{
				{"user_id", 1},
				{"name", 1},
			},
			Options: &options.IndexOptions{
				Unique: &t,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := client.Database(viper.GetString(config.AppName))

	opts := options.Update().SetUpsert(true)
	_, err := db.Collection("counters").UpdateOne(
		ctx,
		bson.D{
			{"_id", "tasks"},
			{"seq", bson.D{{"$exists", false}}},
		},
		bson.D{{"$inc", bson.D{{"seq", 1}}}},
		opts,
	)
	if err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			panic(err)
		}
	}

	return mongodb.CreateIndexes(ctx, db, indexes)
}
