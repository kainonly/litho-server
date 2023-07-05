package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Email       string               `bson:"email" json:"email"`
	Password    string               `bson:"password" json:"-"`
	Roles       []primitive.ObjectID `bson:"roles" json:"-"`
	Name        string               `bson:"name" json:"name"`
	Avatar      string               `bson:"avatar" json:"avatar"`
	BackupEmail string               `bson:"backup_email" json:"backup_email"`
	Lark        UserLark             `bson:"lark" json:"lark"`
	Sessions    int64                `bson:"sessions" json:"sessions"`
	History     UserHistory          `bson:"history" json:"history"`
	Status      bool                 `bson:"status" json:"status"`
	CreateTime  time.Time            `bson:"create_time" json:"create_time"`
	UpdateTime  time.Time            `bson:"update_time" json:"update_time"`
}

type UserLark struct {
	AccessToken      string `bson:"access_token" json:"access_token"`
	TokenType        string `bson:"token_type" json:"token_type"`
	ExpiresIn        uint64 `bson:"expires_in" json:"expires_in"`
	Name             string `bson:"name" json:"name"`
	EnName           string `bson:"en_name" json:"en_name"`
	AvatarUrl        string `bson:"avatar_url" json:"avatar_url"`
	AvatarThumb      string `bson:"avatar_thumb" json:"avatar_thumb"`
	AvatarMiddle     string `bson:"avatar_middle" json:"avatar_middle"`
	AvatarBig        string `bson:"avatar_big" json:"avatar_big"`
	OpenId           string `bson:"open_id" json:"open_id"`
	UnionId          string `bson:"union_id" json:"union_id"`
	Email            string `bson:"email" json:"email"`
	EnterpriseEmail  string `bson:"enterprise_email" json:"enterprise_email"`
	UserId           string `bson:"user_id" json:"user_id"`
	Mobile           string `bson:"mobile" json:"mobile"`
	TenantKey        string `bson:"tenant_key" json:"tenant_key"`
	RefreshExpiresIn uint64 `bson:"refresh_expires_in" json:"refresh_expires_in"`
	RefreshToken     string `bson:"refresh_token" json:"refresh_token"`
	Sid              string `bson:"sid" json:"sid"`
}

type UserHistory struct {
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	ClientIP  string    `bson:"client_ip" json:"client_ip"`
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
