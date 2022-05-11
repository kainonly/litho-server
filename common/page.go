package common

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Page struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 父节点
	Parent interface{} `bson:"parent" json:"parent"`

	// 名称
	Name string `bson:"name" json:"name"`

	// 字体图标
	Icon string `bson:"icon,omitempty" json:"icon,omitempty"`

	// 种类
	Kind string `bson:"kind" json:"kind"`

	// 形式
	Manifest string `bson:"manifest,omitempty" json:"manifest,omitempty"`

	// Schema 定义
	Schema *Schema `bson:"schema,omitempty" json:"schema,omitempty"`

	// 数据源
	Source *Source `bson:"source,omitempty" json:"source,omitempty"`

	// 排序
	Sort int64 `bson:"sort" json:"sort"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

type Schema struct {
	// 命名
	Key string `bson:"key" json:"key"`

	// 字段
	Fields SchemaFields `bson:"fields" json:"fields"`

	// 搜索规则
	Rules []interface{} `bson:"rules,omitempty" json:"rules,omitempty"`

	// 启用事务补偿
	Event *bool `bson:"event,omitempty" json:"event,omitempty"`
}

type SchemaFields map[string]*Field

type Field struct {
	// 显示名称
	Label string `bson:"label" json:"label"`

	// 字段类型
	Type string `bson:"type" json:"type"`

	// 字段描述
	Description string `bson:"description,omitempty" json:"description,omitempty"`

	// 字段提示
	Placeholder string `bson:"placeholder,omitempty" json:"placeholder,omitempty"`

	// 默认值
	Default interface{} `bson:"default,omitempty" json:"default,omitempty"`

	// 是否必须
	Required *bool `bson:"required,omitempty" json:"required,omitempty"`

	// 隐藏字段
	Hide *bool `bson:"hide,omitempty" json:"hide,omitempty"`

	// 只读
	Readonly *bool `bson:"readonly,omitempty" json:"readonly,omitempty"`

	// 排序
	Sort int64 `bson:"sort" json:"sort"`

	// 配置
	Option *FieldOption `bson:"option,omitempty" json:"option,omitempty"`
}

type FieldOption struct {
	// 最大值
	Max int64 `bson:"max,omitempty" json:"max,omitempty"`

	// 最小值
	Min int64 `bson:"min,omitempty" json:"min,omitempty"`

	// 保留小数
	Decimal int64 `bson:"decimal,omitempty" json:"decimal,omitempty"`

	// 包含时间
	Time *bool `bson:"time,omitempty" json:"time,omitempty"`

	// 枚举数值
	Values []FieldValue `bson:"values,omitempty" json:"values,omitempty"`

	// 引用数据源
	Reference string `bson:"reference,omitempty" json:"reference,omitempty"`

	// 引用目标
	Target string `bson:"target,omitempty" json:"target,omitempty"`

	// 是否多选
	Multiple *bool `bson:"multiple,omitempty" json:"multiple,omitempty"`
}

type FieldValue struct {
	// 名称
	Label string `bson:"label" json:"label"`

	// 数值
	Value interface{} `bson:"value" json:"value"`
}

type Source struct {
	// 布局
	Layout string `bson:"layout" json:"layout"`

	// 图表
	Panels []Panel `bson:"panels" json:"panels"`
}

type Panel struct {
	// 模式
	Query string `bson:"query" json:"query"`

	// 映射
	Mappings map[string]string `bson:"mappings" json:"mappings"`

	// 样式
	Style map[string]interface{} `bson:"style,omitempty" json:"style,omitempty"`
}
