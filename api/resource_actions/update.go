package resource_actions
package resource_actions

import (
	"context"
	"time"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID         string `json:"id" vd:"required"`
	ResourceID string `json:"resource_id" vd:"required"`
	Active     *bool  `json:"active" vd:"required"`
	Name       string `json:"name" vd:"required"`
	Code       string `json:"code" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.ResourceActionsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.Db.Model(model.ResourceAction{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(common.M{
			`updated_at`:  time.Now(),
			`resource_id`: dto.ResourceID,
			`active`:      *dto.Active,
			`name`:        dto.Name,
			`code`:        dto.Code,
		}).Error; err != nil {
		return
	}
	return
}
