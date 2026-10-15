package common

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
)

type A []any

func (x *A) Scan(value any) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x A) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`[]`), nil
	}
	return sonic.Marshal(x)
}

type AInt64 []int64

func (x *AInt64) Scan(value any) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x AInt64) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`[]`), nil
	}
	return sonic.Marshal(x)
}

type AString []string

func (x *AString) Scan(value any) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x AString) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`[]`), nil
	}
	return sonic.Marshal(x)
}

type M map[string]any

func (x *M) Scan(value any) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x M) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`{}`), nil
	}
	return sonic.Marshal(x)
}
