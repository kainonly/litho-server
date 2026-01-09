package roles

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) Delete(ctx context.Context, c *app.RequestContext) {
	var dto common.DeleteDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.RolesX.Delete(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]any{"ok": true})
}

func (x *Service) Delete(ctx context.Context, user *common.IAMUser, dto common.DeleteDto) (err error) {
	return
}
