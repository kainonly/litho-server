package model

import (
	"time"
)

type LoginLog struct {
	Timestamp time.Time     `bson:"timestamp"`
	Metadata  LoginMetadata `bson:"metadata"`
	Data      LoginData     `bson:"data"`
}

type LoginMetadata struct {
	UserID  string `bson:"user_id"`
	Email   string `bson:"email"`
	TokenId string `bson:"token_id"`
	Ip      string `bson:"ip"`
	Channel string `bson:"channel" json:"channel"`
}

type LoginData struct {
	Country   string `bson:"country" json:"country"`
	Province  string `bson:"province" json:"province"`
	City      string `bson:"city" json:"city"`
	Isp       string `bson:"isp" json:"isp"`
	UserAgent string `bson:"user_agent"`
}

func NewLoginLog(metadata LoginMetadata, data LoginData) *LoginLog {
	return &LoginLog{
		Timestamp: time.Now(),
		Metadata:  metadata,
		Data:      data,
	}
}
