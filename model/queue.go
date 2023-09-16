package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Queue struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Project     primitive.ObjectID `bson:"project" json:"project"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Subjects    []string           `bson:"subjects" json:"subjects"`
	MaxMsgs     int64              `bson:"max_msgs" json:"max_msgs"`
	MaxBytes    int64              `bson:"max_bytes" json:"max_bytes"`
	MaxAge      time.Duration      `bson:"max_age" json:"max_age"`
	CreateTime  time.Time          `bson:"create_time" json:"create_time" farker:"-"`
	UpdateTime  time.Time          `bson:"update_time" json:"update_time" farker:"-"`
}

func SetQueues(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "queues"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("queue", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "queues", option); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "queues"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}
	return
}
