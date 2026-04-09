package orders

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) FindById(ctx context.Context, c *app.RequestContext) {
	var dto common.FindByIdDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.OrdersX.FindById(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type FindByIdResult struct {
	ID           string     `json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	DepartmentID string     `json:"department_id"`
	UserID       string     `json:"user_id"`
	No           string     `json:"no"`
	Amount       float64    `json:"amount"`
	Status       int16      `json:"status"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	Remark       string     `json:"remark"`
	PaidAt       *time.Time `json:"paid_at"`
	ClosedAt     *time.Time `json:"closed_at"`
}

func (x *Service) FindById(ctx context.Context, user *common.IAMUser, dto common.FindByIdDto) (result FindByIdResult, err error) {
	do := x.Db.Model(model.Order{}).WithContext(ctx)
	ctx = common.SetPipe(ctx, common.NewFindByIdPipe())
	if err = dto.Take(ctx, do, &result); err != nil {
		return
	}
	return
}
