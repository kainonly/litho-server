package routes

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type UpdateDto struct {
	ID     string `json:"id" vd:"required"`
	Nav    string `json:"nav" vd:"required"`
	Pid    string `json:"pid"`
	Type   *int16 `json:"type" vd:"required"`
	Name   string `json:"name" vd:"required"`
	Icon   string `json:"icon"`
	Link   string `json:"link"`
	Status *bool  `json:"status" vd:"required"`
}

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

	common.SetTracking(c, Key, "update", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	data := M{
		"update_time": time.Now(),
		"nav":         dto.Nav,
		"pid":         "0",
		"type":        dto.Type,
		"name":        dto.Name,
		"icon":        "",
		"link":        "",
		"status":      *dto.Status,
	}

	if dto.Pid != "" {
		data["pid"] = dto.Pid
	}
	if dto.Icon != "" {
		data["icon"] = dto.Icon
	}
	if dto.Link != "" {
		data["link"] = dto.Link
	}

	return x.Db.Transaction(func(tx *gorm.DB) (errX error) {
		if errX = tx.Model(model.Route{}).WithContext(ctx).
			Where(`pid = ?`, dto.ID).
			Update(`nav`, dto.Nav).
			Error; errX != nil {
			return
		}
		if errX = tx.Model(model.Route{}).WithContext(ctx).
			Where("id = ?", dto.ID).
			Updates(data).Error; errX != nil {
			return
		}

		return x.RefreshCache(ctx)
	})
}
