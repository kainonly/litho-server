package users

import (
	"context"
	"server/model"
	"time"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type SetStatusesDto struct {
	IDs    []string `json:"ids" vd:"required,min=1"`
	Status *bool    `json:"status" vd:"required"`
}

const ISetStatuses = "批量启用/禁用"

func (x *Controller) SetStatuses(ctx context.Context, c *app.RequestContext) {
	var dto SetStatusesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.UsersX.SetStatuses(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) SetStatuses(ctx context.Context, user *common.IAMUser, dto SetStatusesDto) (err error) {
	updates := common.M{
		`update_time`: time.Now(),
		`status`:      *dto.Status,
	}
	return x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id in (?)`, dto.IDs).
		Updates(updates).Error
}
