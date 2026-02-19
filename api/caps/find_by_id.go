package caps

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
	data, err := x.CapsX.FindById(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type FindByIdResult struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func (x *Service) FindById(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (data *FindByIdResult, err error) {
	var result = new(FindByIdResult)
	if err = x.Db.WithContext(ctx).
		Model(model.Cap{}).
		Where(`id = ?`, dto.ID).
		Scan(result).Error; err != nil {
		return
	}
	data = result
	return
}
