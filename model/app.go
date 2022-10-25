package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type App struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 应用名称
	Name string `bson:"name" json:"name"`

	// Access Key ID
	AppId string `bson:"key" json:"key"`

	// Secret Access Key
	Secret string `bson:"secret" json:"secret"`

	// 后端入口
	Entry []string `bson:"entry" json:"entry"`

	// 有效时间
	ExpireTime time.Time `bson:"expire_time" json:"expire_time"`

	// 状态
	Status bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"create_time"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"update_time"`
}
