package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Department struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 父节点
	Parent *primitive.ObjectID `bson:"parent" json:"parent"`

	// 名称
	Name string `bson:"name" json:"name"`

	// 描述
	Description string `bson:"description" json:"description"`

	// 排序
	Sort int64 `bson:"sort" json:"sort"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

func NewDepartment(name string) *Department {
	return &Department{
		Name:       name,
		Sort:       0,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
