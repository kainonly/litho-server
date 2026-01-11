package permissions

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto common.FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, utils.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	data, err := x.PermissionsX.Find(ctx, user, dto)
	if err != nil {
		c.JSON(500, utils.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto common.FindDto) (data []model.Permission, err error) {
	return
}
