package model

import "go.mongodb.org/mongo-driver/bson"

type Schema struct {
	Key         string  `bson:"key" json:"key"`
	Label       string  `bson:"label" json:"label"`
	Kind        string  `bson:"kind" json:"kind"`
	Description string  `bson:"description,omitempty" json:"description"`
	System      bool    `bson:"system,omitempty" json:"system"`
	Fields      []Field `bson:"fields,omitempty" json:"fields"`
}

type Field struct {
	Key         string      `bson:"key" json:"key"`
	Label       string      `bson:"label" json:"label"`
	Type        string      `bson:"type" json:"type"`
	Description string      `bson:"description,omitempty" json:"description"`
	Default     string      `bson:"default,omitempty" json:"default,omitempty"`
	Unique      bool        `bson:"unique,omitempty" json:"unique,omitempty"`
	Required    bool        `bson:"required,omitempty" json:"required,omitempty"`
	Private     bool        `bson:"private,omitempty" json:"private,omitempty"`
	System      bool        `bson:"system,omitempty" json:"system,omitempty"`
	Option      FieldOption `bson:"option,omitempty" json:"option,omitempty"`
}

type FieldOption struct {
	// 数字类型
	Max interface{} `bson:"max,omitempty" json:"max,omitempty"`
	Min interface{} `bson:"min,omitempty" json:"min,omitempty"`
	// 枚举类型
	Values   bson.D `bson:"values,omitempty" json:"values,omitempty"`
	Multiple *bool  `bson:"multiple,omitempty" json:"multiple,omitempty"`
	// 引用类型
	Mode   string `bson:"mode,omitempty" json:"mode,omitempty"`
	Target string `bson:"target,omitempty" json:"target,omitempty"`
	To     string `bson:"to,omitempty" json:"to,omitempty"`
}
