package orders

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
	Status      int16     `json:"status"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Remark      string    `json:"remark"`
}

const IUpdate = "更新"

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.OrdersX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) error {
	updates := common.M{
		`updated_at`:   time.Now(),
		`status`:       dto.Status,
		`scheduled_at`: dto.ScheduledAt,
		`remark`:       dto.Remark,
	}
	return x.Db.Model(model.Order{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(updates).Error
}
