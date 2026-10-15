package common

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
	"github.com/kainonly/collector/v3/transfer"
	"github.com/kainonly/go/captcha"
	"github.com/kainonly/go/cipher"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/locker"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
	"resty.dev/v3"
)

type Inject struct {
	V        *Values
	Db       *gorm.DB
	RDb      *redis.Client
	MGo      *mongo.Client
	MDb      *mongo.Database
	Nc       *nats.Conn
	Js       jetstream.JetStream
	Cipher   *cipher.Cipher
	Captcha  *captcha.Captcha
	Locker   *locker.Locker
	Transfer *transfer.Transfer
}

// 通用错误定义

var ErrAuthenticationExpired = help.E(4003, "登录账号认证过期，请重新登录")
var ErrLoginNotExists = help.E(4001, "登录账号不存在或被冻结")
var ErrLoginMaxFailures = help.E(4002, "登录账户失败超出最大次数")
var ErrLoginInvalid = help.E(4001, "登录账户不存在或密码不正确")
var ErrSession = help.E(4005, "登录账户会话建立失败")
var ErrTotpInvalid = help.E(4001, "登录口令验证码无效")
var ErrSmsInvalid = help.E(4001, "登录账户短信验证码无效")
var ErrSmsNotExists = help.E(4001, "登录账号不存在或被冻结")
var ErrCodeFrequently = help.E(4004, "您的验证码请求频繁，请稍后再试")

type Fn func(do *gorm.DB) *gorm.DB

type IAMUser struct {
	ID       string       `json:"id"`
	RoleID   string       `json:"role_id" gorm:"-"`
	OrgIDs   string       `json:"org_ids" gorm:"-"`
	Status   bool         `json:"status"`
	Strategy *IAMStrategy `json:"strategy" gorm:"-"`
	Ip       string       `json:"-"`
}

// Can 验证用户是否具备定义细粒化策略
func (x *IAMUser) Can(permissions ...string) error {
	exists := make(map[string]bool)
	for _, permission := range permissions {
		exists[permission] = true
	}
	for _, permission := range x.Strategy.Permissions {
		if exists[permission.(string)] {
			return nil
		}
	}
	return help.E(0, `操作失败，该用户不具备执行权限`)
}

type IAMStrategy struct {
	Navs        A `json:"navs" vd:"required"`
	Routes      A `json:"routes" vd:"required"`
	Permissions A `json:"permissions" vd:"required"`
}

func (x *IAMStrategy) Scan(value any) error {
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

// IpClient 公网 IP 地理信息查询客户端
// 匿名调用存在限额，生产环境可按需追加 token 查询参数
var IpClient = resty.New().
	SetBaseURL(`https://ipinfo.io`).
	SetRetryCount(3)
