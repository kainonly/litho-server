package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Email       string               `bson:"email" json:"email"`
	Roles       []primitive.ObjectID `bson:"permissions" json:"-"`
	Password    string               `bson:"password" json:"-"`
	Name        string               `bson:"name" json:"name"`
	Avatar      string               `bson:"avatar" json:"avatar"`
	BackupEmail string               `bson:"backup_email" json:"backup_email"`
	Sessions    int64                `bson:"sessions" json:"sessions"`
	Last        UserLast             `bson:"last" json:"last"`
	Status      bool                 `bson:"status" json:"status"`
	CreateTime  time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime  time.Time            `bson:"update_time" json:"update_time"`
}

type UserLast struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Ip        string    `bson:"ip" json:"ip"`
	Country   string    `bson:"country" json:"country"`
	Province  string    `bson:"province" json:"province"`
	City      string    `bson:"city" json:"city"`
	Isp       string    `bson:"isp" json:"isp"`
}

func NewUser(email string, password string) *User {
	return &User{
		Email:      email,
		Password:   password,
		Roles:      []primitive.ObjectID{},
		Status:     true,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
