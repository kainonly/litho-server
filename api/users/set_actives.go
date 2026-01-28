package users

import (
	"context"
	"server/model"
	"time"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type SetActivesDto struct {
	IDs    []string `json:"ids" vd:"required,min=1"`
	Active *bool    `json:"active" vd:"required"`
}

func (x *Controller) SetActives(ctx context.Context, c *app.RequestContext) {
	var dto SetActivesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.UsersX.SetActives(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) SetActives(ctx context.Context, user *common.IAMUser, dto SetActivesDto) (err error) {
	updates := common.M{
		`update_time`: time.Now(),
		`active`:      *dto.Active,
	}
	return x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id in (?)`, dto.IDs).
		Updates(updates).Error
}
