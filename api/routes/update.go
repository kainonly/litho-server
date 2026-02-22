package routes

import (
	"context"
	"time"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID     string `json:"id" vd:"required"`
	Active *bool  `json:"active" vd:"required"`
	Pid    string `json:"pid"`
	Name   string `json:"name" vd:"required"`
	Type   *int16 `json:"type"`
	Icon   string `json:"icon"`
	Link   string `json:"link"`
}

const IUpdate = "更新"

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RoutesX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	updates := common.M{
		`updated_at`: time.Now(),
		`active`:     *dto.Active,
		`name`:       dto.Name,
	}
	if dto.Pid != "" {
		updates[`pid`] = dto.Pid
	}
	if dto.Type != nil {
		updates[`type`] = *dto.Type
	}
	if dto.Icon != "" {
		updates[`icon`] = dto.Icon
	}
	if dto.Link != "" {
		updates[`link`] = dto.Link
	}

	if err = x.Db.Model(model.Route{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(updates).Error; err != nil {
		return
	}
	return
}
