package index

import (
	"server/api/sessions"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/kainonly/go/passport"
	"go.uber.org/fx"
)

var Provides = fx.Provide(
	func(i *Service) *Controller { return &Controller{IndexX: i} },
	func(i common.Inject, p *passport.Passport, s *sessions.Service) *Service {
		return &Service{Inject: &i, Passport: p, SessionsX: s}
	},
)

type Controller struct {
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
