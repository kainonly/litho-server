package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Role struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 权限标识
	Key string `bson:"key" json:"key"`

	// 父节点
	Parent interface{} `bson:"parent" json:"parent"`

	// 名称
	Name string `bson:"name" json:"name"`

	// 授权页面
	Pages []primitive.ObjectID `bson:"pages" json:"pages"`

	// 只读权限
	Readonly []primitive.ObjectID `bson:"readonly" json:"readonly"`

	// 状态
	Status *bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

func NewRole(key string, name string) *Role {
	return &Role{
		Key:        key,
		Parent:     nil,
		Name:       name,
		Pages:      []primitive.ObjectID{},
		Readonly:   []primitive.ObjectID{},
		Status:     Bool(true),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

func (x *Role) SetParent(v interface{}) *Role {
	x.Parent = v
	return x
}

func (x *Role) SetPages(v []primitive.ObjectID) *Role {
	x.Pages = v
	return x
}

func (x *Role) SetReadonly(v []primitive.ObjectID) *Role {
	x.Readonly = v
	return x
}
