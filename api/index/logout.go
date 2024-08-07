package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
)

func (x *Controller) Logout(ctx context.Context, c *app.RequestContext) {
	claims := common.Claims(c)
	x.IndexX.Logout(ctx, claims.UserId)
	common.ClearAccessToken(c)
	c.Status(204)
}

func (x *Service) Logout(ctx context.Context, userId string) {
	x.Sessions.Remove(ctx, userId)
}
