package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCreateCollectionForProjects(t *testing.T) {
	ctx := context.TODO()
	option := options.CreateCollection().
		SetValidator(bson.D{
			{"$jsonSchema", bson.D{
				{"title", "projects"},
				{"required", bson.A{"_id", "name", "namespace", "secret", "entry", "expire", "status", "create_time", "update_time"}},
				{"properties", bson.D{
					{"_id", bson.M{"bsonType": "objectId"}},
					{"name", bson.M{"bsonType": "string"}},
					{"namespace", bson.M{"bsonType": "string"}},
					{"secret", bson.M{"bsonType": "string"}},
					{"entry", bson.M{"bsonType": "array"}},
					{"expire", bson.M{"bsonType": "number"}},
					{"status", bson.M{"bsonType": "bool"}},
					{"create_time", bson.M{"bsonType": "date"}},
					{"update_time", bson.M{"bsonType": "date"}},
				}},
				{"additionalProperties", false},
			}},
		})
	if err := db.CreateCollection(ctx, "projects", option); err != nil {
		t.Error(err)
	}
	r, err := db.Collection("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"namespace", 1}},
		Options: options.Index().SetUnique(true).SetName("idx_namespace"),
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestCreateProjects(t *testing.T) {
	var err error
	_, err = db.Collection("projects").InsertOne(
		context.TODO(),
		model.NewProject("默认示例", "default"),
	)
	assert.NoError(t, err)
}
