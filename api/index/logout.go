package index

import (
	"context"
	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

func (x *Controller) Logout(ctx context.Context, c *app.RequestContext) {
	user := common.GetIAM(c)
	x.IndexX.Logout(ctx, user)

	x.IndexX.ClearAccessToken(c)
	c.JSON(200, help.Ok())
}

func (x *Service) Logout(ctx context.Context, user *common.IAMUser) {
	x.SessionsX.Kick(ctx, user, user.ID)
}
