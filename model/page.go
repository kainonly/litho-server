package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"time"
)

type Page struct {
	ID        uint64    `json:"id"`
	Parent    uint64    `gorm:"default:0;not null;comment:父节点" json:"parent"`
	Name      string    `gorm:"type:varchar;not null;comment:名称" json:"name"`
	Icon      string    `gorm:"type:varchar;comment:字体图标" json:"icon,omitempty"`
	Kind      string    `gorm:"type:varchar;default:'default';not null;comment:种类" json:"kind"`
	Manifest  string    `gorm:"type:varchar;default:'default';not null;comment:形式" json:"manifest,omitempty"`
	Schema    Schema    `gorm:"type:jsonb;comment:模型，数据集时存在" json:"schema,omitempty"`
	Source    Source    `gorm:"type:jsonb;comment:数据源，数据聚合时存在" json:"source,omitempty"`
	Manual    Manual    `gorm:"type:jsonb;comment:自定义，自定义种类时存在" json:"manual,omitempty"`
	Sort      int64     `gorm:"default:0;comment:排序" json:"sort"`
	Status    bool      `gorm:"default:true;not null;comment:状态" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Schema struct {
	// 命名
	Key string `json:"key"`

	// 字段
	Fields []SchemaField `json:"fields"`

	// 显隐规则
	Rules []SchemaRule `json:"rules,omitempty"`

	// 启用事务补偿
	Event *bool `json:"event,omitempty"`

	// 启用详情
	Detail *bool `json:"detail,omitempty"`
}

func (x *Schema) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return sonic.Unmarshal(bytes, x)
}

func (x Schema) Value() (driver.Value, error) {
	return sonic.MarshalString(x)
}

type SchemaField struct {
	// 命名
	Key string `json:"key"`

	// 显示名称
	Label string `json:"label"`

	// 字段类型
	Type string `json:"type"`

	// 字段描述
	Description string `json:"description,omitempty"`

	// 字段提示
	Placeholder string `json:"placeholder,omitempty"`

	// 默认值
	Default interface{} `json:"default,omitempty"`

	// 关键词
	Keyword bool `json:"keyword,omitempty"`

	// 是否必须
	Required bool `json:"required,omitempty"`

	// 隐藏字段
	Hide bool `json:"hide,omitempty"`

	// 只读
	Readonly bool `json:"readonly,omitempty"`

	// 投影
	Projection int64 `json:"projection,omitempty"`

	// 排序
	Sort int64 `json:"sort"`

	// 配置
	Option *SchemaFieldOption `json:"option,omitempty"`
}

type SchemaFieldOption struct {
	// 最大值
	Max int64 `json:"max,omitempty"`

	// 最小值
	Min int64 `json:"min,omitempty"`

	// 保留小数
	Decimal int64 `json:"decimal,omitempty"`

	// 包含时间
	Time bool `json:"time,omitempty"`

	// 枚举数值
	Values []Value `json:"values,omitempty"`

	// 引用类型，模型
	Reference string `json:"reference,omitempty"`

	// 引用类型，目标字段
	Target string `json:"target,omitempty"`

	// 多选
	Multiple bool `json:"multiple,omitempty"`

	// 组件标识
	Component string `json:"component,omitempty"`
}

type Value struct {
	// 名称
	Label string `json:"label"`

	// 数值
	Value interface{} `json:"value"`
}

type SchemaRule struct {
	// 逻辑
	Logic string `json:"logic"`

	// 条件
	Conditions []*SchemaRuleCondition `json:"conditions"`

	// 显示字段
	Keys []string `json:"keys"`
}

type SchemaRuleCondition struct {
	// 字段
	Key string `json:"key"`

	// 操作符
	Operate string `json:"operate"`

	// 数值
	Value interface{} `json:"value"`
}

type Source struct {
	// 布局
	Layout string `json:"layout"`

	// 图表
	Panels []*Panel `json:"panels"`
}

func (x *Source) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return sonic.Unmarshal(bytes, x)
}

func (x Source) Value() (driver.Value, error) {
	return sonic.MarshalString(x)
}

type Panel struct {
	// 模式
	Query string `json:"query"`

	// 映射
	Mappings map[string]string `json:"mappings"`

	// 样式
	Style map[string]interface{} `json:"style,omitempty"`
}

type Manual struct {
	// 页面标识，自定义页面接入命名
	Scope string `json:"scope"`

	// 权限细粒化
	Policies map[string]string `json:"policies"`
}

func (x *Manual) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return sonic.Unmarshal(bytes, x)
}

func (x Manual) Value() (driver.Value, error) {
	return sonic.MarshalString(x)
}
