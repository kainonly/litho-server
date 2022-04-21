package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Media struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 类型
	Type string `bson:"type" json:"type"`

	// 名称
	Name string `bson:"name" json:"name"`

	// 媒体地址
	Url string `bson:"url" json:"url"`

	// 数据参数
	Params map[string]string `bson:"params" json:"params"`

	// 标记
	Labels []string `bson:"labels" json:"labels"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}
