package schedulers

import (
	"context"
	"fmt"
	"server/common"
	"server/model"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.SchedulesX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID       string `json:"id"`
	Status   *bool  `json:"status"`
	TeamID   string `json:"team_id"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Secret   string `json:"secret"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(model.Scheduler{}).WithContext(ctx)

	if dto.Q != "" {
		do = do.Where(`name like ?`, fmt.Sprintf(`%%%s%%`, dto.Q))
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	results = make([]*FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}

	return
}
