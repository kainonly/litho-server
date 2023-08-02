package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Category struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Type       string             `bson:"type" json:"type"`
	Name       string             `bson:"name" json:"name"`
	Sort       int64              `bson:"sort" json:"sort"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}
