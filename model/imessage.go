package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Imessage struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Topic       string               `bson:"topic" json:"topic"`
	Description string               `bson:"description" json:"description"`
	Projects    []primitive.ObjectID `bson:"projects" json:"projects"`
	CreateTime  time.Time            `bson:"create_time" json:"create_time" farker:"-"`
	UpdateTime  time.Time            `bson:"update_time" json:"update_time" farker:"-"`
}
