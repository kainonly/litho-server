package roles

import (
	"context"
	"time"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type SetStrategyDto struct {
	ID       string              `json:"id" vd:"required"`
	Strategy common.RoleStrategy `json:"strategy" vd:"required"`
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

	c.JSON(200, help.Ok())
}

func (x *Service) SetStrategy(ctx context.Context, user *common.IAMUser, dto SetStrategyDto) (err error) {
	if err = x.Db.Model(model.Role{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(common.M{
			`updated_at`: time.Now(),
			`strategy`:   dto.Strategy,
		}).Error; err != nil {
		return
	}
	if err = x.RefreshCache(ctx); err != nil {
		return
	}
	return
}
