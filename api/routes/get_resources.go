package routes

import (
	"context"
	"litho-api/common"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) GetResources(ctx context.Context, c *app.RequestContext) {
	if err := x.RoutesX.CheckAccess(ctx, common.GetIAM(c)); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, x.V.Resources.List())
}
