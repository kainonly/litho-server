package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type Array []string

func (x *Array) Scan(input interface{}) error {
	text, ok := input.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to format String value:", input))
	}
	*x = strings.Split(text, ",")
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
	b, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(b, x)
}

func (x JSONObject) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	b, err := jsoniter.Marshal(x)
	return string(b), err
}
