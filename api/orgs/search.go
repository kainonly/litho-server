package orgs

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type SearchDto struct {
	common.SearchDto

	Type int16 `query:"type"`
}

func (x *Controller) Search(ctx context.Context, c *app.RequestContext) {
	var dto SearchDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	data, err := x.OrgsX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type SearchResult struct {
	ID   string `json:"id"`
	Type *int16 `json:"type"`
	Name string `json:"name"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []SearchResult, err error) {
	do := x.Db.Model(model.Org{}).WithContext(ctx)
	if dto.Type != 0 {
		do = do.Where("type = ?", dto.Type)
	}
	if dto.Q != "" {
		do = do.Where("name like ?", dto.GetKeyword())
	}

	results = make([]SearchResult, 0)
	ctx = common.SetPipe(ctx, common.NewSearchPipe(`id`, `type`, `name`).SkipAsync())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
