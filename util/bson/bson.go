package bson

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DocToStruct(d primitive.D, result any) error {
	b, err := bson.Marshal(d)
	if err != nil {
		return err
	}
	err = bson.Unmarshal(b, result)
	if err != nil {
		return err
	}

	return nil
}
