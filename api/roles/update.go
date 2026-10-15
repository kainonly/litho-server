package roles

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID          string `json:"id" vd:"required"`
	Name        string `json:"name" vd:"required"`
	Description string `json:"description"`
	Status      *bool  `json:"status" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RolesX.Update(ctx, user, dto); err != nil {
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

	updates := M{
		"update_time": time.Now(),
		"name":        dto.Name,
		"description": dto.Description,
		"status":      *dto.Status,
	}
	if err = x.Db.Model(model.Role{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(updates).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
