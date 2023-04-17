package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AccelerationTask struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Source     string             `bson:"source" json:"source"`
	Target     string             `bson:"target" json:"target"`
	Status     bool               `bson:"status" json:"status"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}

func NewAccelerationTask(source string, target string) *AccelerationTask {
	return &AccelerationTask{
		Source:     source,
		Target:     target,
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
