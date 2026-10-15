package index

import (
	"context"
	"errors"
	"litho-api/api/sessions"
	"litho-api/api/tencent"
	"litho-api/common"
	"litho-api/model"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/goforj/wire"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/locker"
	"github.com/kainonly/go/passport"
	"gorm.io/gorm"
	"resty.dev/v3"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	V *common.Values

	IndexX *Service
}

type Service struct {
	*common.Inject

	Passport  *passport.Passport
	SessionsX *sessions.Service
	TencentX  *tencent.Service
}

type M = map[string]any

type LoginResult struct {
	*model.User `json:"-"`

	AccessToken string `json:"access_token"`
}

func (x *Service) QueryLogin(ctx context.Context, handleFunc common.Fn) (result *LoginResult, err error) {
	result = new(LoginResult)
	do := x.Db.Model(model.User{}).WithContext(ctx).
		Where(`status = ?`, true)
	if handleFunc != nil {
		do = handleFunc(do)
	}

	if err = do.Take(&result.User).Error; err != nil {
		return
	}

	if err = x.Locker.Check(ctx, result.Phone, 5); err != nil {
		switch {
		case errors.Is(err, locker.ErrNotExists):
			err = nil
			break
		case errors.Is(err, locker.ErrLocked):
			err = common.ErrLoginMaxFailures
			return
		default:
			return
		}
	}
	return
}

func (x *Service) UpdateLogin(ctx context.Context, userId string) (err error) {
	update := M{
		"sessions":  gorm.Expr("sessions + ?", 1),
		"last_time": time.Now(),
	}
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, userId).
		Updates(update).Error; err != nil {
		return
	}
	return
}

type LoginAuditDto struct {
	User *LoginResult `bson:"-"`

	Ts        time.Time      `bson:"ts"`
	Meta      LoginAuditMeta `bson:"meta"`
	Ip        string         `bson:"ip"`
	PrivateIP bool           `bson:"private_ip"`
	Country   string         `bson:"country"`
	Region    string         `bson:"region"`
	City      string         `bson:"city"`
	Org       string         `bson:"org"`
	Agent     string         `bson:"agent"`
}

type LoginAuditMeta struct {
	Plat string `bson:"plat"`
	Mode int    `bson:"mode"` // 登录方式：1 账号，2 短信
	UID  int    `bson:"uid"`
}

type IPData struct {
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
	Org     string `json:"org"`
}

func (x *Controller) LoginAudit(c *app.RequestContext, user *LoginResult, mode int) {
	dto := LoginAuditDto{
		User: user,
		Ts:   time.Now(),
		Meta: LoginAuditMeta{
			Plat: "system",
			Mode: mode,
		},
		Ip:    string(c.GetHeader(x.V.Ip)),
		Agent: string(c.GetHeader("User-Agent")),
	}
	go func() {
		_ = x.IndexX.LoginAuditPush(dto)
	}()
}

func (x *Service) LoginAuditPush(dto LoginAuditDto) (err error) {
	ip := strings.TrimSpace(strings.Split(dto.Ip, ",")[0])
	parsedIP := net.ParseIP(ip)
	dto.PrivateIP = parsedIP == nil || common.IsPrivateIP(parsedIP)
	if !dto.PrivateIP {
		var response *resty.Response
		response, requestErr := common.IpClient.R().Get(ip)
		if requestErr == nil {
			var data IPData
			if decodeErr := sonic.Unmarshal(response.Bytes(), &data); decodeErr == nil {
				dto.Country = data.Country
				dto.Region = data.Region
				dto.City = data.City
				dto.Org = data.Org
			}
		}
	}
	if dto.Meta.UID, err = strconv.Atoi(dto.User.ID); err != nil {
		return
	}
	return x.Transfer.Send("login_audits", dto)
}

func (x *Service) CreateAccessToken(ctx context.Context, userId string) (ts string, err error) {
	jti := help.Uuid7()
	claims := passport.NewClaims(userId, time.Hour*24*7).SetJTI(jti)
	if ts, err = x.Passport.Create(claims); err != nil {
		return
	}
	if status := x.SessionsX.Set(ctx, userId, jti); status != "OK" {
		err = common.ErrSession
		return
	}
	x.Locker.Delete(ctx, userId)
	return
}

func (x *Service) SetAccessToken(c *app.RequestContext, ts string) {
	c.SetCookie("TOKEN", ts, 0, "/", "",
		protocol.CookieSameSiteLaxMode, true, true)
}

func (x *Service) ClearAccessToken(c *app.RequestContext) {
	c.SetCookie("TOKEN", "", 0, "/", "",
		protocol.CookieSameSiteLaxMode, true, true)
}
