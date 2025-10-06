package common

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
)

type A []any

func (x *A) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x A) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`[]`), nil
	}
	return sonic.Marshal(x)
}

type M map[string]any

func (x *M) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x M) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`{}`), nil
	}
	return sonic.Marshal(x)
}
