package users

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto common.FindByIdDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	data, err := x.UsersX.FindById(ctx, user, dto)
	if err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(200, data)
}

func (x *Service) FindById(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (data model.User, err error) {
	ctx = common.SetPipe(ctx, common.NewFindByIdPipe().Omit("password"))
	err = dto.Take(ctx, x.Db.Model(&model.User{}), &data)
	return
}
