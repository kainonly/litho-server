package routes

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type SetStrategyDto struct {
	ID       string   `json:"id" vd:"required"`
	Strategy common.M `json:"strategy" vd:"required"`
}

func (x *Controller) SetStrategy(ctx context.Context, c *app.RequestContext) {
	var dto SetStrategyDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RoutesX.SetStrategy(ctx, user, dto); err != nil {
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
	strategy, err := x.V.Resources.NormalizeStrategy(dto.Strategy)
	if err != nil {
		return help.E(0, err.Error())
	}

	if err = x.Db.Model(model.Route{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Update(`strategy`, common.M(strategy)).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
