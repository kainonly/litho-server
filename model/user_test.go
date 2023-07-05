package model_test

import (
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

func TestCreateCollectionForUsers(t *testing.T) {
	ctx := context.TODO()
	b, err := os.ReadFile("./user.json")
	assert.NoError(t, err)
	var jsonSchema bson.D
	err = bson.UnmarshalExtJSON(b, true, &jsonSchema)
	assert.NoError(t, err)

	n, err := x.Db.ListCollectionNames(ctx, bson.M{"name": "users"})
	assert.NoError(t, err)
	if len(n) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		err = x.Db.CreateCollection(ctx, "users", option)
		assert.NoError(t, err)
	} else {
		err = x.Db.RunCommand(ctx, bson.D{
			{"collMod", "users"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err()
		assert.NoError(t, err)
	}

	index := mongo.IndexModel{
		Keys: bson.D{{"email", 1}},
		Options: options.Index().
			SetUnique(true).
			SetName("idx_email"),
	}
	r, err := x.Db.Collection("users").Indexes().CreateOne(ctx, index)
	assert.NoError(t, err)
	t.Log(r)
}

func TestCreateUser(t *testing.T) {
	hash, err := argon2id.CreateHash("pass@VAN1234", argon2id.DefaultParams)
	assert.NoError(t, err)
	_, err = x.Db.Collection("users").InsertMany(
		context.TODO(),
		[]interface{}{
			model.NewUser("kainonly@qq.com", hash),
		},
	)
	assert.NoError(t, err)
}
