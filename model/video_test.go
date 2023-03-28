package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

func TestCreateVideos(t *testing.T) {
	ctx := context.TODO()
	b, err := os.ReadFile("./video.json")
	assert.NoError(t, err)
	var jsonSchema bson.D
	err = bson.UnmarshalExtJSON(b, true, &jsonSchema)
	assert.NoError(t, err)

	n, err := db.ListCollectionNames(ctx, bson.M{"name": "videos"})
	assert.NoError(t, err)
	if len(n) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		err = db.CreateCollection(ctx, "videos", option)
		assert.NoError(t, err)
	} else {
		err = db.RunCommand(ctx, bson.D{
			{"collMod", "videos"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err()
		assert.NoError(t, err)
	}
}

func TestCreateVideoTags(t *testing.T) {
	ctx := context.TODO()
	b, err := os.ReadFile("./video_tag.json")
	assert.NoError(t, err)
	var jsonSchema bson.D
	err = bson.UnmarshalExtJSON(b, true, &jsonSchema)
	assert.NoError(t, err)

	n, err := db.ListCollectionNames(ctx, bson.M{"name": "video_tags"})
	assert.NoError(t, err)
	if len(n) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		err = db.CreateCollection(ctx, "video_tags", option)
		assert.NoError(t, err)
	} else {
		err = db.RunCommand(ctx, bson.D{
			{"collMod", "video_tags"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err()
		assert.NoError(t, err)
	}
}
