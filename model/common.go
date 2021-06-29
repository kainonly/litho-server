package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type Array []interface{}

func (c *Array) Scan(input interface{}) error {
	b, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(b, c)
}

func (c Array) Value() (driver.Value, error) {
	b, err := jsoniter.Marshal(c)
	return string(b), err
}

type Object map[string]interface{}

func (c *Object) Scan(input interface{}) error {
	b, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(b, c)
}

func (c Object) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	b, err := jsoniter.Marshal(c)
	return string(b), err
}
