package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"time"
)

type Acl struct {
	ID         uint64    `json:"id"`
	Status     *bool     `gorm:"default:true" json:"status"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
	Key        string    `gorm:"varchar(20);not null;unique;comment:访问控制索引" json:"key"`
	Name       string    `gorm:"varchar(20);not null;comment:访问控制名称" json:"name"`
	Acts       Acts      `gorm:"type:jsonb;default:'[]';comment:访问控制单元" json:"acts"`
}

type Acts [2][]Act

type Act struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

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
