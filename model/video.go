package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Video struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name       string               `bson:"name" json:"name"`
	Url        string               `bson:"url" json:"url"`
	Categories []primitive.ObjectID `bson:"categories" json:"categories"`
	CreateTime time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime time.Time            `bson:"update_time" json:"update_time"`
}

func SetupVideo(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "videos"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("video", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "videos", option); err != nil {
			return
		}
		index := []mongo.IndexModel{
			{
				Keys:    bson.D{{"name", 1}},
				Options: options.Index().SetName("idx_name"),
			},
			{
				Keys: bson.D{{"url", 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_url"),
			},
		}
		if _, err = db.Collection("videos").
			Indexes().
			CreateMany(ctx, index); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "videos"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}
	return
}
