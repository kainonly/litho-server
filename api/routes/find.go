package routes

import (
	"context"
	"fmt"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto

	Nav string `query:"nav"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.RoutesX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type FindResult struct {
	ID     string `json:"id"`
	Nav    string `json:"nav"`
	Type   int16  `json:"type"`
	Pid    string `json:"pid"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Status bool   `json:"status"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (results []FindResult, err error) {
	do := x.Db.Model(&model.Route{}).WithContext(ctx)
	if dto.Nav != "" {
		do = do.Where("nav = ?", dto.Nav)
	}
	if dto.Q != "" {
		do = do.Where("name like ?", fmt.Sprintf(`%%%s%%`, dto.Q))
	}

	results = make([]FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe().SkipSort())
	if err = dto.Find(ctx, do.Order(`type`).Order(`sort`), &results); err != nil {
		return
	}
	return
}
