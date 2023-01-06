package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	Email       string               `bson:"email" json:"email"`
	Roles       []primitive.ObjectID `bson:"roles" json:"-"`
	Password    string               `bson:"password" json:"-"`
	Name        string               `bson:"name" json:"name"`
	Avatar      string               `bson:"avatar" json:"avatar"`
	BackupEmail string               `bson:"backup_email" json:"backup_email"`
	Feishu      FeishuUserData       `bson:"feishu" json:"feishu"`
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

// 参数详情
// https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/access_token/create
type FeishuUserData struct {
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
