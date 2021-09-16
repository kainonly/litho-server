package model

import (
	"database/sql/driver"
	jsoniter "github.com/json-iterator/go"
)

type Array []interface{}

func (x *Array) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), x)
}

func (x Array) Value() (driver.Value, error) {
	return jsoniter.Marshal(x)
}

type Object map[string]interface{}

func (x *Object) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), x)
}

func (x Object) Value() (driver.Value, error) {
	return jsoniter.Marshal(x)
}

func True() *bool {
	value := true
	return &value
}

func False() *bool {
	return new(bool)
}

type Schema struct {
	ID      int64   `json:"id"`
	Key     string  `gorm:"type:varchar;not null;unique" json:"key"`
	Kind    string  `gorm:"type:varchar;not null" json:"kind"`
	Columns Columns `gorm:"type:jsonb;default:'{}'" json:"columns"`
	System  *bool   `gorm:"default:false" json:"system"`
}

type Columns map[string]Column

func (x *Columns) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), x)
}

func (x Columns) Value() (driver.Value, error) {
	return jsoniter.Marshal(x)
}

type Column struct {
	Label    string   `json:"label"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Default  string   `json:"default,omitempty"`
	Unique   bool     `json:"unique,omitempty"`
	Require  bool     `json:"require,omitempty"`
	Relation Relation `json:"relation,omitempty"`
	Private  bool     `json:"private,omitempty"`
	System   bool     `json:"system,omitempty"`
}

type Relation struct {
	Mode       string `json:"mode,omitempty"`
	Target     string `json:"target,omitempty"`
	References string `json:"references,omitempty"`
}

type Resource struct {
	ID     int64  `json:"id"`
	Name   string `gorm:"type:varchar;not null" json:"name"`
	Path   string `gorm:"type:varchar;not null;unique" json:"path"`
	Parent string `gorm:"type:varchar;default:'root'" json:"parent"`
	Router Router `gorm:"type:jsonb;default:'{}'" json:"router"`
	Nav    *bool  `gorm:"default:false" json:"nav"`
	Icon   string `gorm:"type:varchar" json:"icon"`
	Sort   int8   `gorm:"default:0" json:"sort"`
}

type Router struct {
	Template string       `json:"template,omitempty"`
	Schema   string       `json:"schema,omitempty"`
	Option   RouterOption `json:"options,omitempty"`
}

type RouterOption struct {
	Fetch   bool         `json:"fetch,omitempty"`
	Columns []ViewColumn `json:"columns,omitempty"`
}

type ViewColumn struct {
	Name string `json:"name"`
}

func (x *Router) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), x)
}

func (x Router) Value() (driver.Value, error) {
	return jsoniter.Marshal(x)
}
