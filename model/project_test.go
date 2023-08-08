package model_test

import (
	"context"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
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

func TestMockProjects(t *testing.T) {
	ctx := context.TODO()
	_, err := x.Db.Collection("projects").DeleteMany(ctx, bson.M{})
	assert.NoError(t, err)
	data := make([]interface{}, 0)
	now := time.Now()
	for i := 0; i < 2000; i++ {
		var project model.Project
		err = faker.FakeData(&project)
		assert.NoError(t, err)
		project.Entry = []string{}
		project.Logo = ""
		project.CreateTime = now
		project.UpdateTime = now
		data = append(data, project)
	}
	_, err = x.Db.Collection("projects").InsertMany(ctx, data)
	assert.NoError(t, err)
}
