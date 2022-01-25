package mock

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMedia(t *testing.T) {
	ctx := context.Background()
	if _, err := Db.Collection("pictures").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"labels": 1},
				Options: options.Index().SetName("idx_labels"),
			},
		},
	); err != nil {
		t.Error(err)
	}
	if _, err := Db.Collection("videos").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"labels": 1},
				Options: options.Index().SetName("idx_labels"),
			},
		},
	); err != nil {
		t.Error(err)
	}
}
