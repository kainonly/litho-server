package model

import (
	"database/sql/driver"
	jsoniter "github.com/json-iterator/go"
)

type I18n struct {
	ZhCn string `json:"zh_cn"`
	EnUs string `json:"en_us"`
}

func (c *I18n) Value() (driver.Value, error) {
	buf, err := jsoniter.Marshal(c)
	return string(buf), err
}

func (c *I18n) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), c)
}
