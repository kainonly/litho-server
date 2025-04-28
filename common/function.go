package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func GetIAM(c *app.RequestContext) (user *IAMUser) {
	v, _ := c.Get("identity")
	return v.(*IAMUser)
}

func SetAccessToken(c *app.RequestContext, ts string) {
	c.SetCookie("TOKEN", ts, -1,
		"/", "", protocol.CookieSameSiteStrictMode, true, true)
}

func ClearAccessToken(c *app.RequestContext) {
	c.SetCookie("TOKEN", "", -1,
		"/", "", protocol.CookieSameSiteStrictMode, true, true)
}

type PushDto struct {
	Key      string // 审计索引
	Action   string // 操作函数
	Status   string // 操作反馈
	Snapshot any    // 参数快照
}

func Sha256hex(s string) string {
	b := sha256.Sum256([]byte(s))
	return hex.EncodeToString(b[:])
}

func Hmacsha256(s, key string) string {
	hashed := hmac.New(sha256.New, []byte(key))
	hashed.Write([]byte(s))
	return string(hashed.Sum(nil))
}

type Int interface {
	int32 | int64
}

func IsInc[T Int](data []T) bool {
	for i := 1; i < len(data); i++ {
		if data[i] < data[i-1] {
			return false
		}
	}
	return true
}

func IsDec[T Int](data []T) bool {
	for i := 1; i < len(data); i++ {
		if data[i] > data[i-1] {
			return false
		}
	}
	return true
}

func ToStatus(v int16) bool {
	if v == 1 {
		return true
	}
	return false
}
