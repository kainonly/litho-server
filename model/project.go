package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Project struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Namespace   string             `bson:"namespace" json:"namespace"`
	Secret      string             `bson:"secret" json:"secret"`
	Entry       []string           `bson:"entry" json:"entry"`
	Expire      int64              `bson:"expire" json:"expire"`
	Status      bool               `bson:"status" json:"status"`
	CreatedTime time.Time          `json:"created_time"`
	UpdatedTime time.Time          `json:"updated_time"`
}

func NewProject(name string, namespace string) *Project {
	return &Project{
		Name:        name,
		Namespace:   namespace,
		Entry:       []string{},
		Expire:      0,
		Status:      true,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
}
