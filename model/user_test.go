package model_test

import (
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestCreateCollectionForUsers(t *testing.T) {
	ctx := context.TODO()
	option := options.CreateCollection().
		SetValidator(bson.D{
			{"$jsonSchema", bson.D{
				{"title", "projects"},
				{"required", bson.A{"_id", "email", "password", "name", "avatar", "permissions", "status", "create_time", "update_time"}},
				{"properties", bson.D{
					{"_id", bson.M{"bsonType": "objectId"}},
					{"email", bson.M{"bsonType": "string"}},
					{"password", bson.M{"bsonType": "string"}},
					{"name", bson.M{"bsonType": "string"}},
					{"avatar", bson.M{"bsonType": "string"}},
					{"permissions", bson.M{"bsonType": "object"}},
					{"sessions", bson.M{"bsonType": []string{"number", "null"}}},
					{"last", bson.M{"bsonType": []string{"string", "null"}}},
					{"status", bson.M{"bsonType": "bool"}},
					{"create_time", bson.M{"bsonType": "date"}},
					{"update_time", bson.M{"bsonType": "date"}},
				}},
				{"additionalProperties", false},
			}},
		})
	if err := db.CreateCollection(ctx, "users", option); err != nil {
		t.Error(err)
	}
	r, err := db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true).SetName("idx_email"),
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}

func TestCreateUser(t *testing.T) {
	hash, err := argon2id.CreateHash("pass@VAN1234", argon2id.DefaultParams)
	assert.NoError(t, err)
	_, err = db.Collection("users").InsertOne(
		context.TODO(),
		model.NewUser("zhangtqx@vip.qq.com", hash),
	)
	assert.NoError(t, err)
}
