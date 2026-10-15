package permissions

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID          string `json:"-"`
	Key         string `json:"key" vd:"required"`
	Description string `json:"description" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.PermissionsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "create", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := &model.Permission{
		ID:          dto.ID,
		Key:         dto.Key,
		Description: dto.Description,
	}
	return x.Db.WithContext(ctx).Create(data).Error
}
