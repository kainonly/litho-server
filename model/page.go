package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Page struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 父节点
	Parent *primitive.ObjectID `bson:"parent" json:"parent"`

	// 名称
	Name string `bson:"name" json:"name"`

	// 字体图标
	Icon string `bson:"icon,omitempty" json:"icon,omitempty"`

	// 种类
	// "default" => 作为数据源，包含数据列表与数据填充等功能
	// "form" => 独立的数据填充页面
	// "dashboard" => 数据分析处理、结果展示功能，如数据汇总、趋势分析
	// "group" => 导航中将其他种类分组显示
	Kind string `bson:"kind" json:"kind"`

	// Schema 定义
	Schema *Schema `bson:"schema,omitempty" json:"schema,omitempty"`

	// 排序
	Sort int64 `bson:"sort" json:"sort"`

	// 状态
	Status *bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

func NewPage(name string, kind string) *Page {
	return &Page{
		Name:       name,
		Kind:       kind,
		Sort:       0,
		Status:     Bool(true),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

func (x *Page) SetID(v primitive.ObjectID) *Page {
	x.ID = v
	return x
}

func (x *Page) SetParent(v *primitive.ObjectID) *Page {
	x.Parent = v
	return x
}

func (x *Page) SetIcon(v string) *Page {
	x.Icon = v
	return x
}

func (x *Page) SetSchema(v *Schema) *Page {
	x.Schema = v
	return x
}
