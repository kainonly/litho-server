package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Video struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Name       string               `bson:"name" json:"name"`
	Url        string               `bson:"url" json:"url"`
	Categories []primitive.ObjectID `bson:"categories" json:"categories"`
	CreateTime time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime time.Time            `bson:"update_time" json:"update_time"`
}
