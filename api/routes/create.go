package routes

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID     string `json:"-"`
	Nav    string `json:"nav" vd:"required"`
	Type   *int16 `json:"type" vd:"required"`
	Name   string `json:"name" vd:"required"`
	Pid    string `json:"pid"`
	Icon   string `json:"icon"`
	Link   string `json:"link"`
	Status *bool  `json:"status" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.RoutesX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "create", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	data := &model.Route{
		ID:     dto.ID,
		Nav:    dto.Nav,
		Type:   dto.Type,
		Pid:    "0",
		Name:   dto.Name,
		Status: dto.Status,
	}

	if dto.Pid != "" {
		data.Pid = dto.Pid
	}
	if dto.Icon != "" {
		data.Icon = dto.Icon
	}
	if dto.Link != "" {
		data.Link = dto.Link
	}

	if err = x.Db.WithContext(ctx).Create(data).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
