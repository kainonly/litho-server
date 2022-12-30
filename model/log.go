package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type LoginLog struct {
	Timestamp time.Time     `bson:"timestamp"`
	Metadata  LoginMetadata `bson:"metadata"`
	Data      LoginData     `bson:"data"`
}

type LoginMetadata struct {
	UserID  string `bson:"user"`
	Email   string `bson:"email"`
	TokenId string `bson:"token_id"`
	Ip      string `bson:"ip"`
}

type LoginData struct {
	Detail    bson.M `bson:"detail"`
	UserAgent string `bson:"user_agent"`
}

func NewLoginLog(metadata LoginMetadata, data LoginData) *LoginLog {
	return &LoginLog{
		Timestamp: time.Now(),
		Metadata:  metadata,
		Data:      data,
	}
}
