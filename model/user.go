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
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"-"`
	Name        string             `bson:"name" json:"name"`
	Avatar      string             `bson:"avatar" json:"avatar"`
	Permissions Permissions        `bson:"permissions" json:"-"`
	Status      bool               `bson:"status" json:"status"`
	CreateTime  time.Time          `bson:"create_time" json:"create_time"`
	UpdateTime  time.Time          `bson:"update_time" json:"update_time"`
}

func NewUser(email string, password string) *User {
	return &User{
		Email:       email,
		Password:    password,
		Permissions: map[string]interface{}{},
		Status:      true,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
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
