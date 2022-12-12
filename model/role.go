package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"time"
)

type Role struct {
	ID          uint64    `json:"id"`
	Name        string    `gorm:"type:varchar;uniqueIndex;not null;comment:名称" json:"name"`
	Description string    `gorm:"type:varchar;comment:描述" json:"description"`
	Pages       RolePages `gorm:"type:jsonb;default:'{}';not null;comment:授权页面" json:"pages"`
	Status      bool      `gorm:"default:true;not null;comment:状态" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RolePages map[string]int64

func (x *RolePages) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return sonic.Unmarshal(bytes, x)
}

func (x RolePages) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	return sonic.MarshalString(x)
}
