package roles

import (
	"context"
	"litho-api/common"
	"litho-api/model"

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
	data, err := x.RolesX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type SearchResult struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *bool  `json:"status"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []SearchResult, err error) {
	do := x.Db.Model(model.Role{}).WithContext(ctx)
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
	}

	results = make([]SearchResult, len(results))
	ctx = common.SetPipe(ctx, common.NewSearchPipe(`id`, `name`, `description`, `status`).SkipAsync())
	if err = dto.Find(ctx, do.Order(`sort`), &results); err != nil {
		return
	}
	return
}
