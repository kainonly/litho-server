package roles

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type UpdateDto struct {
	ID          string  `json:"id" vd:"required"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Sort        *int16  `json:"sort,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, utils.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.RolesX.Update(ctx, user, dto); err != nil {
		c.JSON(500, utils.H{"error": err.Error()})
		return
	}
	c.JSON(200, utils.H{"ok": true})
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	return
}
