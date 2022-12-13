package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"time"
)

type User struct {
	ID          uint64      `json:"id"`
	Email       string      `gorm:"type:varchar;not null;uniqueIndex;comment:电子邮件" json:"email"`
	Password    string      `gorm:"type:varchar;not null;comment:密码" json:"-"`
	Name        string      `gorm:"type:varchar;default:'';not null;comment:称呼" json:"name"`
	Avatar      string      `gorm:"type:varchar;default:'';not null;comment:头像" json:"avatar"`
	Permissions Permissions `gorm:"type:jsonb;default:'{}';not null;comment:授权" json:"-"`
	Status      bool        `gorm:"default:true;not null;comment:状态" json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Permissions map[string]interface{}

func (x *Permissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return sonic.Unmarshal(bytes, x)
}

func (x Permissions) Value() (driver.Value, error) {
	if len(x) == 0 {
		return nil, nil
	}
	return sonic.MarshalString(x)
}
