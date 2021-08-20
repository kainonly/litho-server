package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type Acl struct {
	ID     uint64 `json:"id"`
	Status *bool  `gorm:"default:true" json:"status"`
	Name   string `gorm:"type:varchar(20);not null;comment:访问控制名称" json:"name"`
	Path   string `gorm:"type:varchar(20);not null;unique;comment:访问路径" json:"path"`
	Acts   Acts   `gorm:"type:json;default:'[]';comment:访问控制单元" json:"acts"`
}

type Acts struct {
	R Act `json:"r"`
	W Act `json:"w"`
}

type Act map[string]string

func (x *Acts) Scan(input interface{}) error {
	data, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(data, x)
}

func (x Acts) Value() (driver.Value, error) {
	data, err := jsoniter.Marshal(x)
	return string(data), err
}
