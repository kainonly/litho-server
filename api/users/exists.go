package users

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) Exists(ctx context.Context, c *app.RequestContext) {
	var dto common.ExistsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.UsersX.Exists(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

func (x *Service) Exists(ctx context.Context, user *common.IAMUser, dto common.ExistsDto) (result common.ExistsResult, err error) {
	do := x.Db.Model(model.User{}).WithContext(ctx)
	ctx = common.SetPipe(ctx, common.NewExistsPipe(`email`))
	return dto.Exists(ctx, do)
}
