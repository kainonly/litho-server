package model

import (
	"time"
)

type ValuesLog struct {
	// 时间
	Time time.Time `bson:"time"`

	// 操作用户
	Uid interface{} `bson:"uid"`

	// 快照
	Snapshot interface{} `bson:"snapshot"`
}
