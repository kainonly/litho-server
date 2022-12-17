package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Project struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Namespace  string             `bson:"namespace" json:"namespace"`
	Secret     string             `bson:"secret" json:"secret"`
	Entry      []string           `bson:"entry" json:"entry"`
	Expire     int64              `bson:"expire" json:"expire"`
	Status     bool               `bson:"status" json:"status"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}

func NewProject(name string, namespace string) *Project {
	return &Project{
		Name:       name,
		Namespace:  namespace,
		Entry:      []string{},
		Expire:     0,
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
