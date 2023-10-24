package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Project struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	Name       string              `bson:"name" json:"name" faker:"first_name"`
	Namespace  string              `bson:"namespace" json:"namespace" faker:"username,unique"`
	Kind       string              `bson:"kind" json:"kind"`
	Expire     *time.Time          `bson:"expire" json:"expire" farker:"-"`
	SecretId   string              `bson:"secret_id" json:"secret_id" faker:"unique"`
	SecretKey  string              `bson:"secret_key" json:"secret_key"`
	Entry      []string            `bson:"entry" json:"entry" farker:"-"`
	Status     bool                `bson:"status" json:"status"`
	Cluster    *primitive.ObjectID `bson:"cluster" json:"cluster"`
	Nats       *ProjectNats        `bson:"nats" json:"nats"`
	Redis      *ProjectRedis       `bson:"redis" json:"redis"`
	CreateTime time.Time           `bson:"create_time" json:"create_time" farker:"-"`
	UpdateTime time.Time           `bson:"update_time" json:"update_time" farker:"-"`
}

type ProjectNats struct {
	Seed string `bson:"seed" json:"seed"`
	Pub  string `bson:"pub" json:"pub"`
}

type ProjectRedis struct {
	Url  string `bson:"url" json:"url"`
	Auth string `bson:"auth" json:"auth"`
}

func NewProject(name string, namespace string) *Project {
	return &Project{
		Name:       name,
		Namespace:  namespace,
		Entry:      []string{},
		Expire:     nil,
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

func SetProjects(ctx context.Context, db *mongo.Database) (err error) {
	var ns []string
	if ns, err = db.ListCollectionNames(ctx, bson.M{"name": "projects"}); err != nil {
		return
	}
	var jsonSchema bson.D
	if err = LoadJsonSchema("project", &jsonSchema); err != nil {
		return
	}
	if len(ns) == 0 {
		option := options.CreateCollection().SetValidator(jsonSchema)
		if err = db.CreateCollection(ctx, "projects", option); err != nil {
			return
		}
		index := []mongo.IndexModel{
			{
				Keys: bson.D{{"namespace", 1}},
				Options: options.Index().
					SetUnique(true).
					SetName("idx_namespace"),
			},
		}
		if _, err = db.Collection("projects").
			Indexes().
			CreateMany(ctx, index); err != nil {
			return
		}
	} else {
		if err = db.RunCommand(ctx, bson.D{
			{"collMod", "projects"},
			{"validator", jsonSchema},
			{"validationLevel", "strict"},
		}).Err(); err != nil {
			return
		}
	}
	return
}
