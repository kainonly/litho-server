package orgs

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID     string `json:"id" vd:"required"`
	Type   *int16 `json:"type" vd:"required"`
	Name   string `json:"name" vd:"required"`
	Active *bool  `json:"active" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.OrgsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.Db.Model(model.Org{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(common.M{
			`updated_at`: time.Now(),
			`active`:     *dto.Active,
			`name`:       dto.Name,
		}).Error; err != nil {
		return
	}
	return
}
