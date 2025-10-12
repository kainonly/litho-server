package schedulers

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID       string `json:"-"`
	TeamID   string `json:"team_id" vd:"required"`
	Name     string `json:"name" vd:"required"`
	Endpoint string `json:"endpoint" vd:"required"`
	Secret   string `json:"secret" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	dto.ID = help.SID()
	user := common.GetIAM(c)
	if err := x.SchedulesX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := &model.Scheduler{
		ID:       dto.ID,
		TeamID:   dto.TeamID,
		Name:     dto.Name,
		Endpoint: dto.Endpoint,
		Secret:   dto.Secret,
	}
	return x.Db.WithContext(ctx).Create(data).Error
}
