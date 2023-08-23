package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Workflow struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Project    primitive.ObjectID `bson:"project" json:"project"`
	Name       string             `bson:"name" json:"name"`
	Kind       string             `bson:"kind" json:"kind"`
	Schedule   *WorkflowSchedule  `bson:"schedule" json:"schedule"`
	CreateTime time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime time.Time          `bson:"update_time" json:"update_time"`
}

type WorkflowSchedule struct {
	ScheduleId primitive.ObjectID    `bson:"schedule_id" json:"schedule_id"`
	Status     bool                  `bson:"status" json:"status"`
	Jobs       []WorkflowScheduleJob `bson:"jobs" json:"jobs"`
}

type WorkflowScheduleJob struct {
	Mode   string `bson:"mode" json:"mode"`
	Spec   string `bson:"spec" json:"spec"`
	Option bson.M `bson:"option" json:"option"`
}
