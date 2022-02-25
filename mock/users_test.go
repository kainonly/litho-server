package mock

import (
	"api/model"
	"context"
	"github.com/weplanx/go/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestUsers(t *testing.T) {
	ctx := context.Background()
	if err := Db.Collection("users").Drop(ctx); err != nil {
		t.Error(err)
	}
	var root model.Role
	if err := Db.Collection("roles").FindOne(ctx, bson.M{
		"name": "超级管理员",
	}).Decode(&root); err != nil {
		t.Error(err)
	}
	hashPwd, _ := password.Create("pass@VAN1234")
	data := []interface{}{
		model.NewUser("kain", hashPwd).
			SetRoles([]primitive.ObjectID{root.ID}),
	}
	if _, err := Db.Collection("users").
		InsertMany(ctx, data); err != nil {
		t.Error(err)
	}
	if _, err := Db.Collection("users").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"username": 1},
				Options: options.Index().SetName("uk_username").SetUnique(true),
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
