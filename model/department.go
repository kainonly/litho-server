package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Department struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 父节点
	Parent interface{} `bson:"parent" json:"parent"`

	// 名称
	Name string `bson:"name" json:"name"`

	// 描述
	Description string `bson:"description" json:"description"`

	// 标记
	Labels []string `bson:"labels" json:"labels"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

func NewDepartment(name string) *Department {
	return &Department{
		Name:       name,
		Labels:     []string{},
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

func (x *Department) SetDescription(v string) *Department {
	x.Description = v
	return x
}

func (x *Department) SetLabel(v string) *Department {
	x.Labels = append(x.Labels, v)
	return x
}
