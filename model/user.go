package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 电子邮件
	Email string `bson:"email" json:"email"`

	// 密码
	Password string `bson:"password" json:"-"`

	// 称呼
	Name string `bson:"name" json:"name"`

	// 头像
	Avatar string `bson:"avatar" json:"avatar"`

	// 授权
	Permissions Permissions `bson:"permissions" json:"-"`

	// 状态
	Status bool `bson:"status" json:"status"`

	// 创建时间
	CreatedTime time.Time `bson:"created_time" json:"created_time"`

	// 创建时间
	UpdatedTime time.Time `bson:"updated_time" json:"updated_time"`
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
