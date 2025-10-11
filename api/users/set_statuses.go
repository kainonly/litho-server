package users

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type SetStatusesDto struct {
	IDs    []string `json:"ids" vd:"required"`
	Status *bool    `json:"status" vd:"required"`
}

func (x *Controller) SetStatuses(ctx context.Context, c *app.RequestContext) {
	var dto SetStatusesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.UsersX.SetStatuses(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) SetStatuses(ctx context.Context, user *common.IAMUser, dto SetStatusesDto) (err error) {
	return x.Db.Transaction(func(tx *gorm.DB) (err error) {
		updates := M{
			`update_time`: time.Now(),
			`status`:      dto.Status,
		}
		if err = tx.Model(model.User{}).WithContext(ctx).
			Where(`id in (?)`, dto.IDs).
			Updates(updates).Error; err != nil {
			return
		}
		return
	})
}
