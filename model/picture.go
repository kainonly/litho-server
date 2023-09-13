package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Picture struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name       string               `bson:"name" json:"name"`
	Url        string               `bson:"url" json:"url"`
	Query      string               `bson:"query" json:"query"`
	Process    PictureProcess       `bson:"process" json:"process"`
	Categories []primitive.ObjectID `bson:"categories" json:"categories"`
	CreateTime time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime time.Time            `bson:"update_time" json:"update_time"`
}

type PictureProcess struct {
	Mode int64              `bson:"mode" json:"mode"`
	Cut  PictureProcessCut  `bson:"cut" json:"cut"`
	Zoom PictureProcessZoom `bson:"zoom" json:"zoom"`
}

type PictureProcessCut struct {
	X int64 `bson:"x" json:"x"`
	Y int64 `bson:"y" json:"y"`
	W int64 `bson:"w" json:"w"`
	H int64 `bson:"h" json:"h"`
}

type PictureProcessZoom struct {
	W int64 `bson:"w" json:"w"`
	H int64 `bson:"h" json:"h"`
}

func SetupPicture(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "pictures"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("picture", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "pictures", option); err != nil {
			return
		}
		index := []mongo.IndexModel{
			{
				Keys:    bson.D{{"name", 1}},
				Options: options.Index().SetName("idx_name"),
			},
		}
		if _, err = db.Collection("pictures").
			Indexes().
			CreateMany(ctx, index); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "pictures"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}
	return
}
