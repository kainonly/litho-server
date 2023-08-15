package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Workflow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name       string             `bson:"name" json:"name"`
	Kind       string             `bson:"kind" json:"kind"`
	Option     interface{}        `bson:"option" json:"option"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}
