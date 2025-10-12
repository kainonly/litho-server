package jobs

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID          string    `json:"id" vd:"required"`
	SchedulerID string    `json:"scheduler_id" vd:"required"`
	Name        string    `json:"name" vd:"required"`
	Schema      *common.M `json:"schema" vd:"required"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.JobsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	data := M{
		"update_time":  time.Now(),
		"scheduler_id": dto.SchedulerID,
		"name":         dto.Name,
		"schema":       dto.Schema,
	}

	if err = x.Db.Model(model.Job{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(data).
		Error; err != nil {
		return
	}

	return
}
