package model_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

func TestCreateAccelerationTasks(t *testing.T) {
	ctx := context.TODO()
	b, err := os.ReadFile("./acceleration_task.json")
	assert.NoError(t, err)
	var jsonSchema bson.D
	err = bson.UnmarshalExtJSON(b, true, &jsonSchema)
	assert.NoError(t, err)

	n, err := db.ListCollectionNames(ctx, bson.M{"name": "acceleration_tasks"})
	assert.NoError(t, err)
	if len(n) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		err = db.CreateCollection(ctx, "acceleration_tasks", option)
		assert.NoError(t, err)
	} else {
		err = db.RunCommand(ctx, bson.D{
			{"collMod", "acceleration_tasks"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err()
		assert.NoError(t, err)
	}
}

func TestMockAccelerationTasks(t *testing.T) {
	keys := []string{
		"/ng-zorro-antd@15.1.0/ng-zorro-antd.compact.min.css",
		"/ng-zorro-antd@15.1.0/code-editor/style/index.min.css",
		"/ng-zorro-antd@15.1.0/resizable/style/index.min.css",
		"/@editorjs/editorjs@2.26.5/dist/editor.min.js",
		"/@editorjs/paragraph@2.9.0/dist/bundle.min.js",
		"/@editorjs/header@2.7.0/dist/bundle.min.js",
		"/@editorjs/delimiter@1.3.0/dist/bundle.min.js",
		"/@editorjs/underline@1.1.0/dist/bundle.min.js",
		"/@editorjs/nested-list@1.3.0/dist/nested-list.min.js",
		"/@editorjs/checklist@1.4.0/dist/bundle.min.js",
		"/@editorjs/table@2.2.1/dist/table.min.js",
		"/@editorjs/quote@2.5.0/dist/bundle.js",
		"/@editorjs/code@2.8.0/dist/bundle.js",
		"/@editorjs/marker@1.3.0/dist/bundle.js",
		"/@editorjs/inline-code@1.4.0/dist/bundle.js",
		"/localforage@1.10.0/dist/localforage.min.js",
		"/cropperjs@1.5.13/dist/cropper.min.js",
		"/cropperjs@1.5.13/dist/cropper.min.css",
		"/monaco-editor@0.36.1/min/vs/editor/editor.main.min.js",
		"/monaco-editor@0.36.1/min/vs/editor/editor.main.min.css",
	}
	data := make([]interface{}, len(keys))
	for i, key := range keys {
		data[i] = model.NewAccelerationTask(
			"https://cdn.jsdelivr.net/npm"+key,
			"npm"+key,
		)
	}
	_, err := db.Collection("acceleration_tasks").InsertMany(context.TODO(), data)
	assert.NoError(t, err)
}
