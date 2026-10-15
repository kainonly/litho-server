package routes

import (
	"context"
	"fmt"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type SearchDto struct {
	common.SearchDto

	Nav  string `query:"nav"`
	Type *int16 `query:"type"`
}

func (x *Controller) Search(ctx context.Context, c *app.RequestContext) {
	var dto SearchDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.RoutesX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type SearchResult struct {
	ID     string `json:"id"`
	Nav    string `json:"nav"`
	Type   int16  `json:"type"`
	Pid    string `json:"pid"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Status bool   `json:"status"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []*SearchResult, err error) {
	do := x.Db.Model(&model.Route{}).WithContext(ctx)
	if dto.Nav != "" {
		do = do.Where("nav = ?", dto.Nav)
	}
	if dto.Type != nil {
		do = do.Where("type = ?", *dto.Type)
	}
	if dto.Q != "" {
		do = do.Where("name like ?", fmt.Sprintf(`%%%s%%`, dto.Q))
	}

	results = make([]*SearchResult, 0)
	ctx = common.SetPipe(ctx, common.NewSearchPipe(`id`, `nav`, `type`, `pid`, `name`, `icon`, `status`).SkipAsync())
	if err = dto.Find(ctx, do.Order(`type`).Order(`sort`), &results); err != nil {
		return
	}
	return
}
