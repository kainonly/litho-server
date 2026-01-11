package routes

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID     string `json:"-"`
	MenuID string `json:"menu_id" vd:"required"`
	Sort   int16  `json:"sort"`
	Active *bool  `json:"active" vd:"required"`
	PID    string `json:"pid"`
	Name   string `json:"name" vd:"required"`
	Type   int16  `json:"type" vd:"required"`
	Icon   string `json:"icon"`
	Link   string `json:"link"`
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

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.Route{
		ID:     dto.ID,
		MenuID: dto.MenuID,
		Sort:   dto.Sort,
		Active: *dto.Active,
		PID:    dto.PID,
		Name:   dto.Name,
		Type:   dto.Type,
		Icon:   dto.Icon,
		Link:   dto.Link,
	}
	if err = x.Db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return
	}
	return
}
