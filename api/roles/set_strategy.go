package roles

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type SetStrategyDto struct {
	ID       string              `json:"id" vd:"required"`
	Strategy *common.IAMStrategy `json:"strategy" vd:"required"`
}

func (x *Controller) SetStrategy(ctx context.Context, c *app.RequestContext) {
	var dto SetStrategyDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RolesX.SetStrategy(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "set_strategy", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) SetStrategy(ctx context.Context, user *common.IAMUser, dto SetStrategyDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	updates := M{
		`update_time`: time.Now(),
		`strategy`:    dto.Strategy,
	}
	if err = x.Db.Model(model.Role{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(updates).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
