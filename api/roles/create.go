package roles

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

type CreateDto struct {
	Name        string `json:"name" vd:"required"`
	Description string `json:"description"`
	Sort        int16  `json:"sort"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, utils.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.RolesX.Create(ctx, user, dto); err != nil {
		c.JSON(500, utils.H{"error": err.Error()})
		return
	}
	c.JSON(201, utils.H{"ok": true})
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	return
}
