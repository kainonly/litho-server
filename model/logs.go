package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LoginLogDto struct {
	Time      time.Time          `bson:"time"`
	V         string             `bson:"v"`
	User      primitive.ObjectID `bson:"user"`
	Username  string             `bson:"username"`
	Email     string             `bson:"email"`
	TokenId   string             `bson:"token_id"`
	Ip        string             `bson:"ip"`
	Detail    bson.M             `bson:"detail"`
	UserAgent string             `bson:"user_agent"`
}

func NewLoginLogV10(data User, jti string, ip string, agent string) *LoginLogDto {
	return &LoginLogDto{
		Time:      time.Now(),
		V:         "v1.0",
		User:      data.ID,
		Username:  data.Username,
		Email:     data.Email,
		TokenId:   jti,
		Ip:        ip,
		UserAgent: agent,
	}
}
