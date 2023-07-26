package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LogsetLogined struct {
	Timestamp time.Time        `bson:"timestamp" json:"timestamp"`
	Metadata  LoginLogMetadata `bson:"metadata" json:"metadata"`
	Country   string           `bson:"country" json:"country"`
	Province  string           `bson:"province" json:"province"`
	City      string           `bson:"city" json:"city"`
	Isp       string           `bson:"isp" json:"isp"`
	UserAgent string           `bson:"user_agent"`
}

func (x *LogsetLogined) SetUserID(v primitive.ObjectID) {
	x.Metadata.UserID = v
}

type LoginLogMetadata struct {
	UserID   primitive.ObjectID `bson:"user_id"`
	ClientIP string             `bson:"client_ip"`
	Channel  string             `bson:"channel" json:"channel"`
}

func (x *LogsetLogined) SetLocation(v map[string]interface{}) {
	x.Country = v["country"].(string)
	x.Province = v["province"].(string)
	x.City = v["city"].(string)
	x.Isp = v["isp"].(string)
}

func NewLogsetLogin(channel string, ip string, useragent string) *LogsetLogined {
	return &LogsetLogined{
		Timestamp: time.Now(),
		Metadata: LoginLogMetadata{
			Channel:  channel,
			ClientIP: ip,
		},
		UserAgent: useragent,
	}
}
