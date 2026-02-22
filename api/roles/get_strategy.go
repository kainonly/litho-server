package roles

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) GetStrategy(ctx context.Context, c *app.RequestContext) {
	var dto common.FindByIdDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.RolesX.GetStrategy(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

func (x *Service) GetStrategy(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (result *common.RoleStrategy, err error) {
	var data model.Role
	if err = x.Db.Model(model.Role{}).WithContext(ctx).
		Select(`strategy`).
		Where(`id = ?`, dto.ID).
		Take(&data).Error; err != nil {
		return
	}
	result = &data.Strategy
	return
}
