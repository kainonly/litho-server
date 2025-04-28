package common

import (
	"database/sql/driver"
	"time"

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

type Action struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Actions []Action

func (x *Actions) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x Actions) Value() (driver.Value, error) {
	return sonic.Marshal(x)
}

type History struct {
	LastTime time.Time `json:"last_time"`
}

func (x *History) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x *History) Value() (driver.Value, error) {
	if x == nil {
		return sonic.Marshal(History{})
	}
	return sonic.Marshal(x)
}
