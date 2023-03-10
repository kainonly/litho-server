package model

import (
	"time"
)

type LoginLog struct {
	Timestamp time.Time        `bson:"timestamp" json:"timestamp"`
	Metadata  LoginLogMetadata `bson:"metadata" json:"metadata"`
	Country   string           `bson:"country" json:"country"`
	Province  string           `bson:"province" json:"province"`
	City      string           `bson:"city" json:"city"`
	Isp       string           `bson:"isp" json:"isp"`
	UserAgent string           `bson:"user_agent"`
}

func (x *LoginLog) SetUserID(v string) {
	x.Metadata.UserID = v
}

type LoginLogMetadata struct {
	UserID   string `bson:"user_id"`
	ClientIP string `bson:"client_ip"`
	Channel  string `bson:"channel" json:"channel"`
}

func (x *LoginLog) SetLocation(v map[string]interface{}) {
	x.Country = v["country"].(string)
	x.Province = v["province"].(string)
	x.City = v["city"].(string)
	x.Isp = v["isp"].(string)
}

func NewLoginLog(channel string, ip string, useragent string) *LoginLog {
	return &LoginLog{
		Timestamp: time.Now(),
		Metadata: LoginLogMetadata{
			Channel:  channel,
			ClientIP: ip,
		},
		UserAgent: useragent,
	}
}
