package routes

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateDto struct {
	ID     string  `json:"id" vd:"required"`
	MenuID *string `json:"menu_id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Type   *int16  `json:"type,omitempty"`
	Icon   *string `json:"icon,omitempty"`
	Link   *string `json:"link,omitempty"`
	PID    *string `json:"pid,omitempty"`
	Sort   *int16  `json:"sort,omitempty"`
	Active *bool   `json:"active,omitempty"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.RoutesX.Update(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]any{"ok": true})
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	return
}
