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
	if _, err := Db.Collection("media").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"type": 1},
				Options: options.Index().SetName("idx_type"),
			},
			{
				Keys:    bson.M{"labels": 1},
				Options: options.Index().SetName("idx_labels"),
			},
		},
	); err != nil {
		t.Error(err)
	}
}
