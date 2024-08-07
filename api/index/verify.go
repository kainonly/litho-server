package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/common"
)

func (x *Controller) Verify(ctx context.Context, c *app.RequestContext) {
	ts := c.Cookie("TOKEN")
	if ts == nil {
		c.JSON(401, M{
			"code":    0,
			"message": common.ErrAuthenticationExpired.Error(),
		})
		return
	}

	if _, err := x.IndexX.Verify(ctx, string(ts)); err != nil {
		common.ClearAccessToken(c)
		c.JSON(401, M{
			"code":    0,
			"message": common.ErrAuthenticationExpired.Error(),
		})
		return
	}

	c.Status(204)
}

func (x *Service) Verify(ctx context.Context, ts string) (claims passport.Claims, err error) {
	if claims, err = x.Passport.Verify(ts); err != nil {
		return
	}
	result := x.Sessions.Verify(ctx, claims.UserId, claims.ID)
	if !result {
		err = common.ErrSessionInconsistent
		return
	}

	// TODO: Check user status

	x.Sessions.Renew(ctx, claims.UserId)

	return
}
