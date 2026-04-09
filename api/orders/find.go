package orders

import (
	"context"
	"server/common"
	"server/model"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
	DepartmentID string     `query:"department_id,omitempty"`
	UserID       string     `query:"user_id,omitempty"`
	Status       *int16     `query:"status,omitempty"`
	StartTime    *time.Time `query:"start_time,omitempty"`
	EndTime      *time.Time `query:"end_time,omitempty"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.OrdersX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID           string     `json:"id"`
	CreatedAt    *time.Time `json:"created_at"`
	DepartmentID string     `json:"department_id"`
	UserID       string     `json:"user_id"`
	No           string     `json:"no"`
	Amount       float64    `json:"amount"`
	Status       int16      `json:"status"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	Remark       string     `json:"remark"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(&model.Order{}).WithContext(ctx)
	if dto.DepartmentID != "" {
		do = do.Where(`department_id = ?`, dto.DepartmentID)
	}
	if dto.UserID != "" {
		do = do.Where(`user_id = ?`, dto.UserID)
	}
	if dto.Status != nil {
		do = do.Where(`status = ?`, *dto.Status)
	}
	if dto.StartTime != nil {
		do = do.Where(`created_at >= ?`, dto.StartTime)
	}
	if dto.EndTime != nil {
		do = do.Where(`created_at < ?`, dto.EndTime)
	}
	if dto.Q != "" {
		do = do.Where(`no like ?`, dto.GetKeyword())
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	results = make([]*FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe().SkipTs())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
