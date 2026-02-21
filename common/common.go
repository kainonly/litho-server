package common

import (
    "database/sql/driver"

    "github.com/bytedance/sonic"
    "github.com/cloudwego/hertz/pkg/app"
    "github.com/kainonly/go/captcha"
    "github.com/kainonly/go/cipher"
    "github.com/kainonly/go/help"
    "github.com/kainonly/go/locker"
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

type Inject struct {
    V       *Values
    Db      *gorm.DB
    RDb     *redis.Client
    Nc      *nats.Conn
    Js      jetstream.JetStream
    Cipher  *cipher.Cipher
    Captcha *captcha.Captcha
    Locker  *locker.Locker
}

var (
    ErrAuthenticationExpired = help.E(0, "身份验证已过期，请重新登录")
    ErrLoginMaxFailures      = help.E(0, "登录尝试次数超过最大限制")
    ErrLoginNotExists        = help.E(0, "账号不存在或已被禁用")
    ErrLoginInvalid          = help.E(0, "账号不存在或密码错误")
    ErrSession               = help.E(0, "会话建立失败")
    ErrSmsInvalid            = help.E(0, "短信验证码无效")
    ErrWMAccess              = help.E(0, `您不具备权限，请联系超级管理员进行管理`)
)

type HandleFunc func(do *gorm.DB) *gorm.DB

type Action struct {
    Label string `json:"label"`
    Value string `json:"value"`
}

type Actions []Action

func (x *Actions) Scan(value interface{}) error {
    return sonic.Unmarshal(value.([]byte), &x)
}

func (x Actions) Value() (driver.Value, error) {
    if len(x) == 0 {
        return []byte(`[]`), nil
    }
    return sonic.Marshal(x)
}

type IAMUser struct {
    ID       string        `json:"id"`
    RoleID   string        `json:"role_id"`
    OrgID    string        `json:"org_id"`
    OrgType  int16         `json:"org_type"`
    Active   bool          `json:"active"`
    Strategy *RoleStrategy `json:"strategy"`
    Ip       string        `json:"-"`
}

type RoleStrategy struct {
    Navs   []string `json:"navs" vd:"required"`
    Routes []string `json:"routes" vd:"required"`
    Caps   []string `json:"caps" vd:"required"`
}

func (x *RoleStrategy) Scan(value interface{}) error {
    return sonic.Unmarshal(value.([]byte), &x)
}

func (x *RoleStrategy) Value() (driver.Value, error) {
    if x == nil {
        return sonic.Marshal(RoleStrategy{
            Navs:   make([]string, 0),
            Routes: make([]string, 0),
            Caps:   make([]string, 0),
        })
    }
    return sonic.Marshal(*x)
}

// 能力标识常量，由 sync-meta 同步到数据库，常量名即数据库中的 key

const (
    WM = "最高权限，拥有系统所有权限"
)

// Can 验证用户是否具备定义能力
func (x *IAMUser) Can(caps ...string) error {
    exists := make(map[string]bool)
    for _, v := range caps {
        exists[v] = true
    }
    for _, v := range x.Strategy.Caps {
        if exists[v] {
            return nil
        }
    }
    return help.E(0, `操作失败，该团队成员不具备执行权限`)
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

type M = map[string]any

type Object map[string]any

func (x *Object) Scan(value interface{}) error {
    return sonic.Unmarshal(value.([]byte), &x)
}

func (x Object) Value() (driver.Value, error) {
    if len(x) == 0 {
        return []byte(`{}`), nil
    }
    return sonic.Marshal(x)
}

type Array []any

func (x *Array) Scan(value interface{}) error {
    return sonic.Unmarshal(value.([]byte), &x)
}

func (x Array) Value() (driver.Value, error) {
    if len(x) == 0 {
        return []byte(`[]`), nil
    }
    return sonic.Marshal(x)
}
