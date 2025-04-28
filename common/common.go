package common

import (
	"database/sql/driver"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/weplanx/go/captcha"
	"github.com/weplanx/go/cipher"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"gorm.io/gorm"
)

type Inject struct {
	V       *Values
	Db      *gorm.DB
	RDb     *redis.Client
	NC      *nats.Conn
	Js      nats.JetStreamContext
	Jv      nats.KeyValue
	Cipher  *cipher.Cipher
	Captcha *captcha.Captcha
	Locker  *locker.Locker
}

// 通用错误定义

var ErrAuthenticationExpired = errors.NewPublic("认证过期，请重新登录")
var ErrLoginNotExists = help.E(0, "登录账号不存在或被冻结")
var ErrLoginMaxFailures = errors.NewPublic("登录失败超出最大次数")
var ErrLoginInvalid = help.E(0, "登录验证无效")
var ErrSession = errors.NewPrivate("会话建立失败")
var ErrSessionInconsistent = errors.NewPublic("会话令牌不一致")
var ErrTotpInvalid = errors.NewPublic("口令验证码无效")
var ErrSmsInvalid = errors.NewPublic("短信验证码无效")
var ErrSmsNotExists = errors.NewPublic("该账户不存在或被冻结")
var ErrEmailInvalid = errors.NewPublic("邮箱验证码无效")
var ErrEmailNotExists = errors.NewPublic("该账户不存在或被冻结")
var ErrCodeFrequently = errors.NewPublic("您的验证码请求频繁，请稍后再试")
var ErrForbidden = help.E(0, `操作失败，账号不具备管理权限`)

type HandleFunc func(do *gorm.DB) *gorm.DB

type IAMUser struct {
	ID           string       `json:"id"`
	Pid          string       `json:"pid"`
	RoleID       string       `json:"role_id"`
	DepartmentID string       `json:"department_id"`
	Status       bool         `json:"status"`
	Strategy     *IAMStrategy `json:"-"`
	Ip           string       `json:"-"`
}

// IsRoot 是否为负责人
func (x *IAMUser) IsRoot() bool {
	return x.Pid == "0"
}

func (x *IAMUser) Can(permissions ...string) (skip bool, err error) {
	// 验证用户是否具备定义细粒化策略，最高优先级
	exists := make(map[string]bool)
	for _, permission := range permissions {
		exists[permission] = true
	}
	for _, permission := range x.Strategy.Permissions {
		if exists[permission.(string)] {
			skip = true // 存在则跳过数据所属验证
			return
		}
	}
	return
}

type IAMStrategy struct {
	Navs        A `json:"navs" vd:"required"`
	Routes      A `json:"routes" vd:"required"`
	Permissions A `json:"permissions" vd:"required"`
}

func (x *IAMStrategy) Scan(value interface{}) error {
	return sonic.Unmarshal(value.([]byte), &x)
}

func (x *IAMStrategy) Value() (driver.Value, error) {
	return sonic.Marshal(x)
}

type IAMRoute struct {
	ID       string `json:"id"`
	Link     string `json:"link"`
	Strategy M      `json:"strategy"`
}

type NzNode struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	IsLeaf bool   `json:"isLeaf"`
}

type NzTreeNodes struct {
	*NzNode
	Children []*NzNode `json:"children"`
}
