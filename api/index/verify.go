package index

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passport"
)

func (x *Controller) Verify(ctx context.Context, c *app.RequestContext) {
	ts := c.Cookie("TOKEN")
	if ts == nil {
		c.JSON(401, M{
			"code":    0,
			"message": `Auth expired. Please re-login`,
		})
		return
	}

	if _, err := x.IndexX.Verify(ctx, string(ts)); err != nil {
		x.IndexX.ClearAccessToken(c)
		c.JSON(401, M{
			"code":    0,
			"message": `Auth expired. Please re-login`,
		})
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Verify(ctx context.Context, ts string) (claims passport.Claims, err error) {
	if claims, err = x.Passport.Verify(ts); err != nil {
		return
	}
	result := x.SessionsX.Verify(ctx, claims.ActiveId, claims.ID)
	if !result {
		err = help.E(0, `Authentication token mismatch`)
		return
	}
	x.SessionsX.Renew(ctx, claims.ActiveId)
	return
}
