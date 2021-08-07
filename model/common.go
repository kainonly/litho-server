package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
	"time"
)

type Common struct {
	ID         uint64    `json:"id"`
	Status     *bool     `gorm:"default:true" json:"status"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`
}

type Array []string

func (x *Array) Scan(input interface{}) error {
	text, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to format String value:", input))
	}
	*x = strings.Split(string(text), ",")
	return nil
}

func (x Array) Value() (driver.Value, error) {
	return strings.Join(x, ","), nil
}

type JSONArray []interface{}

func (x *JSONArray) Scan(input interface{}) error {
	b, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(b, x)
}

func (x JSONArray) Value() (driver.Value, error) {
	b, err := jsoniter.Marshal(x)
	return string(b), err
}

type JSONObject map[string]interface{}

func (x *JSONObject) Scan(input interface{}) error {
	data, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(data, x)
}

func (x JSONObject) Value() (driver.Value, error) {
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
