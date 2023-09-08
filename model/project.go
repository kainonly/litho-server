package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Project struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name       string             `bson:"name" json:"name" faker:"first_name"`
	Namespace  string             `bson:"namespace" json:"namespace" faker:"username,unique"`
	Kind       string             `bson:"kind" json:"kind"`
	SecretId   string             `bson:"secret_id" json:"secret_id" faker:"unique"`
	SecretKey  string             `bson:"secret_key" json:"secret_key"`
	Entry      []string           `bson:"entry" json:"entry" farker:"-"`
	Expire     *time.Time         `bson:"expire" json:"expire" farker:"-"`
	Status     bool               `bson:"status" json:"status"`
	CreateTime time.Time          `bson:"create_time" json:"create_time" farker:"-"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time" farker:"-"`
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
