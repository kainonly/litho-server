package types

import (
	"database/sql/driver"
	jsoniter "github.com/json-iterator/go"
)

type JSON jsoniter.RawMessage

func (c *JSON) Value() (driver.Value, error) {
	buf, err := jsoniter.Marshal(c)
	return string(buf), err
}

func (c *JSON) Scan(input interface{}) error {
	return jsoniter.Unmarshal(input.([]byte), c)
}
