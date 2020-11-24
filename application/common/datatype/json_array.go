package datatype

import (
	"database/sql/driver"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type JSONArray []interface{}

func (c *JSONArray) Scan(input interface{}) error {
	bs, ok := input.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", input))
	}
	return jsoniter.Unmarshal(bs, c)
}

func (c JSONArray) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	bs, err := jsoniter.Marshal(c)
	return string(bs), err
}
