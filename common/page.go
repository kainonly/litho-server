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

	// 状态
	Status *bool `bson:"status" json:"status"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"create_time"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"update_time"`
}

type Schema struct {
	// 命名
	Key string `bson:"key" json:"key"`

	// 字段
	Fields []*SchemaField `bson:"fields" json:"fields"`

	// 显隐规则
	Rules []*SchemaRule `bson:"rules,omitempty" json:"rules,omitempty"`

	// 启用事务补偿
	Event *bool `bson:"event,omitempty" json:"event,omitempty"`

	// 启用详情
	Detail *bool `bson:"detail,omitempty" json:"detail,omitempty"`
}

type SchemaField struct {
	// 命名
	Key string `bson:"key" json:"key"`

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

	// 关键词
	Keyword *bool `bson:"keyword,omitempty" json:"keyword,omitempty"`

	// 是否必须
	Required *bool `bson:"required,omitempty" json:"required,omitempty"`

	// 隐藏字段
	Hide *bool `bson:"hide,omitempty" json:"hide,omitempty"`

	// 只读
	Readonly *bool `bson:"readonly,omitempty" json:"readonly,omitempty"`

	// 投影
	Projection *int64 `bson:"projection,omitempty" json:"projection,omitempty"`

	// 排序
	Sort *int64 `bson:"sort" json:"sort"`

	// 配置
	Option *SchemaFieldOption `bson:"option,omitempty" json:"option,omitempty"`
}

type SchemaFieldOption struct {
	// 最大值
	Max int64 `bson:"max,omitempty" json:"max,omitempty"`

	// 最小值
	Min int64 `bson:"min,omitempty" json:"min,omitempty"`

	// 保留小数
	Decimal int64 `bson:"decimal,omitempty" json:"decimal,omitempty"`

	// 包含时间
	Time *bool `bson:"time,omitempty" json:"time,omitempty"`

	// 枚举数值
	Values []Value `bson:"values,omitempty" json:"values,omitempty"`

	// 引用类型，模型
	Reference string `bson:"reference,omitempty" json:"reference,omitempty"`

	// 引用类型，目标字段
	Target string `bson:"target,omitempty" json:"target,omitempty"`

	// 多选
	Multiple *bool `bson:"multiple,omitempty" json:"multiple,omitempty"`
}

type Value struct {
	// 名称
	Label string `bson:"label" json:"label"`

	// 数值
	Value interface{} `bson:"value" json:"value"`
}

type SchemaRule struct {
	// 逻辑
	Logic string `bson:"logic" json:"logic"`

	// 条件
	Conditions []*SchemaRuleCondition `bson:"conditions" json:"conditions"`

	// 显示字段
	Keys []string `bson:"keys" json:"keys"`
}

type SchemaRuleCondition struct {
	// 字段
	Field string `bson:"field" json:"field"`

	// 操作符
	Operate string `bson:"operate" json:"operate"`

	// 数值
	Value interface{} `bson:"value" json:"value"`
}

type Source struct {
	// 布局
	Layout string `bson:"layout" json:"layout"`

	// 图表
	Panels []*Panel `bson:"panels" json:"panels"`
}

type Panel struct {
	// 模式
	Query string `bson:"query" json:"query"`

	// 映射
	Mappings map[string]string `bson:"mappings" json:"mappings"`

	// 样式
	Style map[string]interface{} `bson:"style,omitempty" json:"style,omitempty"`
}
