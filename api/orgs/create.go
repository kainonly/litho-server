package orgs

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
)

type CreateDto struct {
	Name string `json:"name" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.OrgsX.Create(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(201, map[string]any{"ok": true})
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	return
}
