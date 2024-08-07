package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/common"
	"time"
)

type RefreshTokenDto struct {
	Code string `json:"code" vd:"required"`
}

func (x *Controller) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var dto RefreshTokenDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := common.Claims(c)
	ts, err := x.IndexX.RefreshToken(ctx, claims, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	common.SetAccessToken(c, ts)
	c.Status(204)
}

func (x *Service) RefreshToken(ctx context.Context, claims passport.Claims, code string) (ts string, err error) {
	if err = x.Captcha.Verify(ctx, claims.UserId, code); err != nil {
		return
	}
	if ts, err = x.Passport.Create(claims.UserId, claims.ID, time.Hour*2); err != nil {
		return
	}
	return
}
