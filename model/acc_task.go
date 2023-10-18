package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccTask struct {
	Kind   string `bson:"kind" json:"kind"`
	Source string `bson:"source" json:"source"`
	Target string `bson:"target" json:"target"`
}

func SetAccTasks(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "acc_tasks"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("acc_task", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "acc_tasks", option); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "acc_tasks"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}

	return
}
