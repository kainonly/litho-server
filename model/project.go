package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Project struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 项目名称
	Name string `bson:"name" json:"name"`

	// 命名空间
	Namespace string `bson:"namespace" json:"namespace"`

	// Secret
	Secret string `bson:"secret" json:"secret"`

	// 后端入口
	Entry []string `bson:"entry" json:"entry"`

	// 有效时间
	ExpireTime *time.Time `bson:"expire_time" json:"expire_time"`

	// 状态
	Status bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"create_time"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"update_time"`
}
