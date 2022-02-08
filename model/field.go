package model

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
		Required: Bool(false),
		Hide:     Bool(false),
		Modified: Bool(true),
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
	x.Required = Bool(true)
	return x
}

func (x *Field) SetHide() *Field {
	x.Hide = Bool(true)
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
