package orgs

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto common.FindByIdDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, utils.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	data, err := x.OrgsX.FindById(ctx, user, dto)
	if err != nil {
		c.JSON(500, utils.H{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (x *Service) FindById(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (data model.Org, err error) {
	return
}
