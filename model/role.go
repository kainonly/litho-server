package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Role struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 名称
	Name string `bson:"name" json:"name"`

	// 描述
	Description string `bson:"description" json:"description"`

	// 授权页面
	Pages map[string]*int64 `bson:"pages" json:"pages"`

	// 标签
	Labels map[string]string `bson:"labels" json:"labels"`

	// 状态
	Status bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

func NewRole(name string) *Role {
	return &Role{
		Name:       name,
		Pages:      map[string]*int64{},
		Labels:     map[string]string{},
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

func (x *Role) SetDescription(v string) *Role {
	x.Description = v
	return x
}
