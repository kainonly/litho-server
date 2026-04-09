package departments

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
	Status *bool  `json:"status" vd:"required"`
}

const IUpdate = "更新"

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.DepartmentsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.Db.Model(model.Department{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(common.M{
			`update_time`: time.Now(),
			`status`:      *dto.Status,
			`type`:        *dto.Type,
			`name`:        dto.Name,
		}).Error; err != nil {
		return
	}
	return
}
