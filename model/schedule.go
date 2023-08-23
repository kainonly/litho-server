package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Schedule struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	NodeId     string             `bson:"node_id" json:"node_id"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}
