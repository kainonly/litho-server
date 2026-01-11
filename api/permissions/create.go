package permissions

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
)

type CreateDto struct {
	Code        string `json:"code" vd:"required"`
	Description string `json:"description" vd:"required"`
	Active      *bool  `json:"active,omitempty"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.PermissionsX.Create(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(201, map[string]any{"ok": true})
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	return
}
