package resource_actions

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID         string `json:"-"`
	ResourceID string `json:"resource_id" vd:"required"`
	Active     *bool  `json:"active" vd:"required"`
	Name       string `json:"name" vd:"required"`
	Code       string `json:"code" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.ResourceActionsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.ResourceAction{
		ID:         dto.ID,
		ResourceID: dto.ResourceID,
		Active:     *dto.Active,
		Name:       dto.Name,
		Code:       dto.Code,
	}
	if err = x.Db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return
	}
	return
}
