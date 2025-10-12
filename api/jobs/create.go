package jobs

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID          string    `json:"-"`
	TeamID      string    `json:"team_id" vd:"required"`
	SchedulerID string    `json:"scheduler_id" vd:"required"`
	Name        string    `json:"name" vd:"required"`
	Schema      *common.M `json:"schema" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	dto.ID = help.SID()
	user := common.GetIAM(c)
	if err := x.JobsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := &model.Job{
		ID:          dto.ID,
		TeamID:      dto.TeamID,
		SchedulerID: dto.SchedulerID,
		Name:        dto.Name,
		Schema:      dto.Schema,
	}
	return x.Db.WithContext(ctx).Create(data).Error
}
