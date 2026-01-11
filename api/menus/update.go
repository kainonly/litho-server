package menus

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
	Sort   int16  `json:"sort"`
	Active *bool  `json:"active" vd:"required"`
	Name   string `json:"name" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.MenusX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.Db.Model(model.Menu{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(utils.H{
			`updated_at`: time.Now(),
			`sort`:       dto.Sort,
			`active`:     *dto.Active,
			`name`:       dto.Name,
		}).Error; err != nil {
		return
	}
	return
}
