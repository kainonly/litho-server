package routes

import (
	"context"
	"time"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID     string `json:"id" vd:"required"`
	MenuID string `json:"menu_id" vd:"required"`
	Active *bool  `json:"active" vd:"required"`
	PID    string `json:"pid"`
	Name   string `json:"name" vd:"required"`
	Type   int16  `json:"type" vd:"required"`
	Icon   string `json:"icon"`
	Link   string `json:"link"`
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

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.Db.Model(model.Route{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(utils.H{
			`updated_at`: time.Now(),
			`menu_id`:    dto.MenuID,
			`active`:     *dto.Active,
			`pid`:        dto.PID,
			`name`:       dto.Name,
			`type`:       dto.Type,
			`icon`:       dto.Icon,
			`link`:       dto.Link,
		}).Error; err != nil {
		return
	}
	return
}
