package common

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func StartDate(date string) time.Time {
	ts, _ := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:00", date))
	return ts
}

func EndDate(date string) time.Time {
	ts, _ := time.Parse(time.DateTime, fmt.Sprintf("%s 00:00:00", date))
	return ts.Add(24 * time.Hour)
}

func GetSegments(c *app.RequestContext) []string {
	v, _ := c.Get("segments")
	return v.([]string)
}

func GetIAM(c *app.RequestContext) *IAMUser {
	v, _ := c.Get("identity")
	return v.(*IAMUser)
}

type Tracking struct {
	RES  string
	ACT  string
	RIDS []string
}

// SetTracking 设置审计追踪
// rid = 0 说明无资源关联，例如 sort，save 等
// rid != 0 用于资源关联，例如产品中的 set_bids 每个 IDs 都是产品的关联，即对应 rids
func SetTracking(c *app.RequestContext, res string, act string, rids ...string) {
	c.Set("tracking", &Tracking{
		RES:  res,
		ACT:  act,
		RIDS: rids,
	})
}

func GetTracking(c *app.RequestContext) *Tracking {
	v, ok := c.Get("tracking")
	if !ok {
		return nil
	}
	return v.(*Tracking)
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

func IsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	return ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}
