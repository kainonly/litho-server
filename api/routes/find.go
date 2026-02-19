package routes

import (
	"context"
	"strconv"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
	Nav string `query:"nav"`
	Pid string `query:"pid"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.RoutesX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID     string `json:"id"`
	MenuID string `json:"menu_id"`
	Active bool   `json:"active"`
	Pid    string `json:"pid"`
	Name   string `json:"name"`
	Type   int16  `json:"type"`
	Icon   string `json:"icon"`
	Link   string `json:"link"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(&model.Route{}).WithContext(ctx)
	if dto.Nav != "" {
		do = do.Where(`nav = ?`, dto.Nav)
	}
	if dto.Pid != "" {
		do = do.Where(`pid = ?`, dto.Pid)
	}
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	results = make([]*FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe())
	if err = dto.Find(ctx, do.Order(`type`).Order(`sort`), &results); err != nil {
		return
	}
	return
}
