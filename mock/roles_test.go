package mock

import (
	"api/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestRoles(t *testing.T) {
	ctx := context.Background()
	if err := Db.Collection("roles").Drop(ctx); err != nil {
		t.Error(err)
	}
	data := []interface{}{
		model.NewRole("*", "超级管理员"),
	}
	if _, err := Db.Collection("roles").
		InsertMany(ctx, data); err != nil {
		t.Error(err)
	}
	if _, err := Db.Collection("roles").Indexes().CreateOne(ctx,
		mongo.IndexModel{
			Keys:    bson.M{"key": 1},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		t.Error(err)
	}
}
