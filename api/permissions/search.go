package permissions

import (
	"context"
	"fmt"
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
	data, err := x.PermissionsX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type SearchResult struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Description string `json:"description"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []*SearchResult, err error) {
	do := x.Db.Model(&model.Permission{}).WithContext(ctx)
	if dto.Q != "" {
		do = do.Where("description like ?", fmt.Sprintf(`%%%s%%`, dto.Q))
	}

	results = make([]*SearchResult, 0)
	ctx = common.SetPipe(ctx, common.NewSearchPipe(`key`, `description`).SkipAsync())
	if err = dto.Find(ctx, do.Order(`key`), &results); err != nil {
		return
	}
	return
}
