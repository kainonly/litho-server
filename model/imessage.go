package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Imessage struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Topic       string               `bson:"topic" json:"topic"`
	Description string               `bson:"description" json:"description"`
	Projects    []primitive.ObjectID `bson:"projects" json:"projects"`
	CreateTime  time.Time            `bson:"create_time" json:"create_time" farker:"-"`
	UpdateTime  time.Time            `bson:"update_time" json:"update_time" farker:"-"`
}

func SetupImessage(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "imessages"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("imessage", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "imessages", option); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "imessages"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}
	return
}
