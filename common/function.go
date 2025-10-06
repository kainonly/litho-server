package common

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func SetAccessToken(c *app.RequestContext, ts string) {
	c.SetCookie("TOKEN", ts, -1,
		"/", "", protocol.CookieSameSiteStrictMode, true, true)
}

func ClearAccessToken(c *app.RequestContext) {
	c.SetCookie("TOKEN", "", -1,
		"/", "", protocol.CookieSameSiteStrictMode, true, true)
}

func GetIAM(c *app.RequestContext) *IAMUser {
	v, _ := c.Get("identity")
	return v.(*IAMUser)
}
