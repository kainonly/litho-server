package model

import (
	"api/common"
	"go.mongodb.org/mongo-driver/bson"
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
		Parent:     nil,
		Kind:       kind,
		Sort:       0,
		Status:     common.BoolToP(true),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}

func (x *Page) SetID(v primitive.ObjectID) *Page {
	x.ID = v
	return x
}

func (x *Page) SetParent(v primitive.ObjectID) *Page {
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

type Schema struct {
	// 集合命名
	Key string `bson:"key" json:"key"`

	// 字段定义
	Fields SchemaFields `bson:"fields" json:"fields"`

	// 规则
	Rules []interface{} `bson:"rules,omitempty" json:"rules,omitempty"`

	// 验证器
	Validator bson.M `bson:"validator,omitempty" json:"validator,omitempty"`
}

type SchemaFields map[string]*Field

func NewSchema(key string, fields SchemaFields) *Schema {
	return &Schema{
		Key:    key,
		Fields: fields,
	}
}

func (x *Schema) SetRules(v []interface{}) *Schema {
	x.Rules = v
	return x
}

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

	// 可编辑
	Modified *bool `bson:"modified,omitempty" json:"modified,omitempty"`

	// 排序
	Sort int64 `bson:"sort" json:"sort"`

	// 规格
	Spec *FieldSpec `bson:"spec" json:"spec"`
}

func NewField(label string, datatype string) *Field {
	return &Field{
		Label:    label,
		Type:     datatype,
		Required: common.BoolToP(false),
		Hide:     common.BoolToP(false),
		Modified: common.BoolToP(true),
		Sort:     0,
	}
}

func (x *Field) SetDescription(v string) *Field {
	x.Description = v
	return x
}

func (x *Field) SetPlaceholder(v string) *Field {
	x.Placeholder = v
	return x
}

func (x *Field) SetDefault(v interface{}) *Field {
	x.Default = v
	return x
}

func (x *Field) SetRequired() *Field {
	x.Required = common.BoolToP(true)
	return x
}

func (x *Field) SetHide() *Field {
	x.Hide = common.BoolToP(true)
	return x
}

func (x *Field) SetModified(v *bool) *Field {
	x.Modified = v
	return x
}

func (x *Field) SetSort(v int64) *Field {
	x.Sort = v
	return x
}

func (x *Field) SetSpec(v *FieldSpec) *Field {
	x.Spec = v
	return x
}

type FieldSpec struct {
	// 最大值
	Max int64 `bson:"max,omitempty" json:"max,omitempty"`

	// 最小值
	Min int64 `bson:"min,omitempty" json:"min,omitempty"`

	// 保留小数
	Decimal int64 `bson:"decimal,omitempty" json:"decimal,omitempty"`

	// 枚举数值
	Values []Value `bson:"values,omitempty" json:"values,omitempty"`

	// 引用数据源
	Reference string `bson:"reference,omitempty" json:"reference,omitempty"`

	// 引用目标
	Target string `bson:"target,omitempty" json:"target,omitempty"`

	// 是否多选
	Multiple *bool `bson:"multiple,omitempty" json:"multiple,omitempty"`
}

type Value struct {
	// 名称
	Label string `bson:"label" json:"label"`

	// 数值
	Value interface{} `bson:"value" json:"value"`
}
