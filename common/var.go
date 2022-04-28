package common

import (
	"github.com/thoas/go-funk"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Var struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Key   string             `json:"key" bson:"key"`
	Value interface{}        `json:"value" bson:"value"`
}

func NewVar(key string, value interface{}) *Var {
	return &Var{
		Key:   key,
		Value: value,
	}
}

var secrets = []string{
	"tencent_secret_key",
	"tencent_pulsar_token",
	"aliyun_access_key_secret",
	"huaweicloud_access_key_secret",
	"feishu_app_secret",
	"feishu_encrypt_key",
	"feishu_verification_token",
	"email_password",
	"openapi_secret",
}

func SecretKey(key string) bool {
	return funk.Contains(secrets, key)
}
