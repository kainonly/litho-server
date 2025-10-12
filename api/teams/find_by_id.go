package teams

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto common.FindByIdDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.TeamsX.FindById(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type FindByIdResult struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (x *Service) FindById(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (result FindByIdResult, err error) {
	do := x.Db.Model(model.User{}).WithContext(ctx)
	ctx = common.SetPipe(ctx, common.NewFindByIdPipe())
	if err = dto.Take(ctx, do, &result); err != nil {
		return
	}

	return
}
