package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

func TestCreateProjects(t *testing.T) {
	ctx := context.TODO()
	b, err := os.ReadFile("./project.json")
	assert.NoError(t, err)
	var jsonSchema bson.D
	err = bson.UnmarshalExtJSON(b, true, &jsonSchema)
	assert.NoError(t, err)

	n, err := db.ListCollectionNames(ctx, bson.M{"name": "projects"})
	assert.NoError(t, err)
	if len(n) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		err = db.CreateCollection(ctx, "projects", option)
		assert.NoError(t, err)
	} else {
		err = db.RunCommand(ctx, bson.D{
			{"collMod", "projects"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err()
		assert.NoError(t, err)
	}

	r, err := db.Collection("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{"namespace", 1}},
		Options: options.Index().SetUnique(true).SetName("idx_namespace"),
	})
	assert.NoError(t, err)
	t.Log(r)
}

func TestMockProject(t *testing.T) {
	var err error
	_, err = db.Collection("projects").InsertOne(
		context.TODO(),
		model.NewProject("默认示例", "default"),
	)
	assert.NoError(t, err)
}
