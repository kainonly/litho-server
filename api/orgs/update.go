package orgs

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateDto struct {
	ID     string  `json:"id" vd:"required"`
	Name   *string `json:"name,omitempty"`
	Active *bool   `json:"active,omitempty"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, utils.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.OrgsX.Update(ctx, user, dto); err != nil {
		c.JSON(500, utils.H{"error": err.Error()})
		return
	}
	c.JSON(200, utils.H{"ok": true})
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	return
}
