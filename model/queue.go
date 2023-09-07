package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Queue struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Project     primitive.ObjectID `bson:"project" json:"project"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Subjects    []string           `bson:"subjects" json:"subjects"`
	MaxMsgs     int64              `bson:"max_msgs" json:"max_msgs"`
	MaxBytes    int64              `bson:"max_bytes" json:"max_bytes"`
	MaxAge      time.Duration      `bson:"max_age" json:"max_age"`
	CreateTime  time.Time          `bson:"create_time" json:"create_time" farker:"-"`
	UpdateTime  time.Time          `bson:"update_time" json:"update_time" farker:"-"`
}
