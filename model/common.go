package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

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
