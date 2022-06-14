package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Schedule struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 名称
	Name string `bson:"name"`

	// 描述
	Description string `bson:"description" json:"description"`

	// 任务
	Jobs []*ScheduleJob `bson:"jobs" json:"jobs"`

	// 状态
	Status *bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

type ScheduleJob struct {
	// 触发模式
	Mode string `json:"mode"`

	// 时间规格
	Spec string `json:"spec"`

	// 配置
	Option map[string]interface{} `json:"option"`
}
