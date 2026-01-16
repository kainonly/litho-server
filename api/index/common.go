package index

import (
	"context"
	"errors"
	"server/api/sessions"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/goforj/wire"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/locker"
	"github.com/kainonly/go/passport"
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
}

func (x *Service) SetAccessToken(c *app.RequestContext, ts string) {
	if x.V.Cors.SameSite == "none" {
		c.SetCookie("TOKEN", ts, -1, "/", "",
			protocol.CookieSameSiteNoneMode, false, false)
		return
	}
	c.SetCookie("TOKEN", ts, -1, "/", "",
		protocol.CookieSameSiteStrictMode, true, true)
}

func (x *Service) ClearAccessToken(c *app.RequestContext) {
	if x.V.Cors.SameSite == "none" {
		c.SetCookie("TOKEN", "", -1, "/", "",
			protocol.CookieSameSiteNoneMode, false, false)
		return
	}
	c.SetCookie("TOKEN", "", -1, "/", "",
		protocol.CookieSameSiteStrictMode, true, true)
}

type LoginResult struct {
	*model.User `json:"-"`

	AccessToken string `json:"access_token"`
}

func (x *Service) QueryLoginUser(ctx context.Context, handleFunc common.HandleFunc) (result *LoginResult, err error) {
	result = new(LoginResult)
	do := x.Db.Model(model.User{}).WithContext(ctx).
		Where(`active = ?`, true)
	if handleFunc != nil {
		do = handleFunc(do)
	}

	if err = do.Take(&result.User).Error; err != nil {
		return
	}

	if err = x.Locker.Check(ctx, result.Phone, 5); err != nil {
		switch {
		case errors.Is(err, locker.ErrNotExists):
			return
		case errors.Is(err, locker.ErrLocked):
			err = common.ErrLoginMaxFailures
			return
		default:
			return
		}
	}
	return
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
