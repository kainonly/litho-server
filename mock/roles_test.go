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
		model.NewRole("超级管理员").
			SetDescription("系统默认设置").
			SetLabel("最高权限"),
	}
	if _, err := Db.Collection("roles").
		InsertMany(ctx, data); err != nil {
		t.Error(err)
	}
	if _, err := Db.Collection("roles").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"name": 1},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.M{"labels": 1},
				Options: options.Index(),
			},
		},
	); err != nil {
		t.Error(err)
	}
}
