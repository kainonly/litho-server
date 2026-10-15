package permissions

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID          string `json:"id" vd:"required"`
	Key         string `json:"key" vd:"required"`
	Description string `json:"description" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.PermissionsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "update", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	updates := M{
		"key":         dto.Key,
		"description": dto.Description,
	}
	return x.Db.Model(model.Permission{}).WithContext(ctx).
		Where("id = ?", dto.ID).
		Updates(updates).Error
}
