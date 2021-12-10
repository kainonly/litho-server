package model

type FieldSpec struct {
	// 最大值
	Max int64 `bson:"max,omitempty" json:"max,omitempty"`

	// 最小值
	Min int64 `bson:"min,omitempty" json:"min,omitempty"`

	// 保留小数
	Decimal int64 `bson:"decimal,omitempty" json:"decimal,omitempty"`

	// 枚举数值
	Values []Enum `bson:"values,omitempty" json:"values,omitempty"`

	// 引用数据源
	Reference string `bson:"reference,omitempty" json:"reference,omitempty"`

	// 引用目标
	Target string `bson:"target,omitempty" json:"target,omitempty"`

	// 是否多选
	Multiple *bool `bson:"multiple,omitempty" json:"multiple,omitempty"`
}
