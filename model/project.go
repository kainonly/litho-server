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

	// 密钥
	Secret string `bson:"secret" json:"secret"`

	// 后端入口
	Entry []string `bson:"entry" json:"entry"`

	// 有效期 TTL
	Expire int64 `bson:"expire" json:"expire"`

	// 状态
	Status bool `bson:"status" json:"status"`

	// 创建时间
	CreatedTime time.Time `json:"created_time"`

	// 创建时间
	UpdatedTime time.Time `json:"updated_time"`
}
