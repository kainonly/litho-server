package roles

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID          string `json:"-"`
	Name        string `json:"name" vd:"required"`
	Description string `json:"description"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	dto.ID = help.SID()
	user := common.GetIAM(c)
	if err := x.RolesX.Create(ctx, user, dto); err != nil {
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

	data := &model.Role{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		Strategy: common.M{
			"routes":      []common.A{},
			"permissions": []common.A{},
		},
	}
	if err = x.Db.WithContext(ctx).
		Create(data).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
