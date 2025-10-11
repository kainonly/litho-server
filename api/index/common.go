package index

import (
	"context"
	"server/api/sessions"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/google/wire"
	"github.com/kainonly/go/help"
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

type M = map[string]any

func (x *Service) SetAccessToken(c *app.RequestContext, ts string) {
	c.SetCookie("TOKEN", ts, -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
}

func (x *Service) ClearAccessToken(c *app.RequestContext) {
	c.SetCookie("TOKEN", "", -1, "/", "", protocol.CookieSameSiteStrictMode, true, true)
}

type LoginResult struct {
	*model.User `json:"-"`

	AccessToken string `json:"access_token"`
}

func (x *Service) QueryLoginUser(ctx context.Context, handleFunc common.HandleFunc) (result *LoginResult, err error) {
	result = new(LoginResult)
	do := x.Db.Model(model.User{}).WithContext(ctx).
		Where(`status = ?`, true)
	if handleFunc != nil {
		do = handleFunc(do)
	}

	if err = do.Take(&result.User).Error; err != nil {
		return
	}

	return
}

func (x *Service) CreateAccessToken(ctx context.Context, userId string) (ts string, err error) {
	jti := help.Uuid()
	claims := passport.NewClaims(userId, time.Hour*24*7).
		SetIssuer(`litho`).
		SetJTI(jti)
	if ts, err = x.Passport.Create(claims); err != nil {
		return
	}
	if status := x.SessionsX.Set(ctx, userId, jti); status != "OK" {
		err = help.E(0, `failed to establish the login account session`)
		return
	}
	x.Locker.Delete(ctx, userId)
	return
}
