package common

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/kainonly/go/captcha"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/locker"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Inject struct {
	fx.In

	V       *Values
	Db      *gorm.DB
	RDb     *redis.Client
	Captcha *captcha.Captcha
	Locker  *locker.Locker
}

// Common error definitions

var ErrAuthenticationExpired = help.E(4003, "Authentication expired, please login again")
var ErrLoginNotExists = help.E(4001, "Account does not exist or is disabled")
var ErrLoginMaxFailures = help.E(4002, "Login attempts exceeded maximum limit")
var ErrLoginInvalid = help.E(4001, "Account does not exist or password is incorrect")
var ErrSession = help.E(4005, "Failed to establish session")
var ErrTotpInvalid = help.E(4001, "Invalid TOTP verification code")
var ErrSmsInvalid = help.E(4001, "Invalid SMS verification code")
var ErrSmsNotExists = help.E(4001, "Account does not exist or is disabled")
var ErrCodeFrequently = help.E(4004, "Verification code requested too frequently, please try again later")

type HandleFunc func(do *gorm.DB) *gorm.DB

func SetAccessToken(c *app.RequestContext, ts string) {
	c.SetCookie("TOKEN", ts, -1,
		"/", "", protocol.CookieSameSiteStrictMode, true, true)
}

func ClearAccessToken(c *app.RequestContext) {
	c.SetCookie("TOKEN", "", -1,
		"/", "", protocol.CookieSameSiteStrictMode, true, true)
}

type IAMUser struct {
	ID     string `json:"id"`
	OrgID  string `json:"org_id"`
	RoleID string `json:"role_id"`
	Active bool   `json:"active"`
	Ip     string `json:"-"`
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

type M map[string]any

func (x *M) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x M) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`{}`), nil
	}
	return sonic.Marshal(x)
}

type A []any

func (x *A) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x A) Value() (driver.Value, error) {
	if len(x) == 0 {
		return []byte(`[]`), nil
	}
	return sonic.Marshal(x)
}
