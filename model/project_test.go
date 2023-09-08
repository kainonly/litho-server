package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

func TestCreateProjects(t *testing.T) {
	ctx := context.TODO()
	b, err := os.ReadFile("./model/project.json")
	assert.NoError(t, err)
	var jsonSchema bson.D
	err = bson.UnmarshalExtJSON(b, true, &jsonSchema)
	assert.NoError(t, err)

	n, err := x.Db.ListCollectionNames(ctx, bson.M{"name": "projects"})
	assert.NoError(t, err)
	if len(n) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		err = x.Db.CreateCollection(ctx, "projects", option)
		assert.NoError(t, err)
	} else {
		err = x.Db.RunCommand(ctx, bson.D{
			{"collMod", "projects"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err()
		assert.NoError(t, err)
	}

	r, err := x.Db.Collection("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"namespace", 1}},
		Options: options.Index().SetUnique(true).SetName("idx_namespace"),
	})
	assert.NoError(t, err)
	t.Log(r)
}
