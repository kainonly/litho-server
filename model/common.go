package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type Datastore struct {
	ID     uint   `json:"id"`
	Key    string `gorm:"type:varchar(50);not null;unique" json:"key"`
	Type   string `gorm:"type:varchar(20);default:collection;not null;comment:类型" json:"type"`
	Schema Schema `gorm:"type:jsonb;default:'{}';comment:模型声明" json:"schema"`
	Lock   *bool  `gorm:"default:false;comment:锁定" json:"lock"`
}

type Schema []Column

func (x *Schema) Scan(input interface{}) error {
	data, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(data, x)
}

func (x Schema) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	data, err := jsoniter.Marshal(x)
	return string(data), err
}

type Column struct {
	// 字段名称
	Key string `json:"key"`

	// 显示名称
	Label string `json:"label"`

	// 数据类型
	Type string `json:"type"`

	// 默认值
	Default string `json:"default,omitempty"`

	// 必须的, not null
	Require *bool `json:"require,omitempty"`

	// 唯一的
	Unique *bool `json:"unique,omitempty"`

	// 最大长度
	Length uint `json:"length,omitempty"`

	// 备注
	Comment string `json:"comment,omitempty"`

	// 隐藏字段
	Hide *bool `json:"hide,omitempty"`
}

type Array []interface{}

func (x *Array) Scan(input interface{}) error {
	b, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(b, x)
}

func (x Array) Value() (driver.Value, error) {
	b, err := jsoniter.Marshal(x)
	return string(b), err
}

type Object map[string]interface{}

func (x *Object) Scan(input interface{}) error {
	data, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(data, x)
}

func (x Object) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	data, err := jsoniter.Marshal(x)
	return string(data), err
}

func True() *bool {
	value := true
	return &value
}

func False() *bool {
	return new(bool)
}
