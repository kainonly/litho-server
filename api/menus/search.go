package menus

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type SearchDto struct {
	common.SearchDto
}

func (x *Controller) Search(ctx context.Context, c *app.RequestContext) {
	var dto SearchDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	results, err := x.MenusX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, results)
}

type SearchResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []*SearchResult, err error) {
	do := x.Db.Model(model.Menu{}).WithContext(ctx)
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
	}

	results = make([]*SearchResult, 0)
	ctx = common.SetPipe(ctx, common.NewSearchPipe())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
