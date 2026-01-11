package users

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type ExistsDto struct {
	common.ExistsDto
}

func (x *Controller) Exists(ctx context.Context, c *app.RequestContext) {
	var dto ExistsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	result, err := x.UsersX.Exists(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, result)
}

func (x *Service) Exists(ctx context.Context, user *common.IAMUser, dto ExistsDto) (result common.ExistsResult, err error) {
	do := x.Db.Model(model.User{}).WithContext(ctx)
	ctx = common.SetPipe(ctx, common.NewExistsPipe(`email`, `phone`))
	return dto.Exists(ctx, do)
}
