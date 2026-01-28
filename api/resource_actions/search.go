package resource_actions

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type SearchDto struct {
	common.SearchDto
	ResourceID string `json:"resource_id" query:"resource_id"`
}

func (x *Controller) Search(ctx context.Context, c *app.RequestContext) {
	var dto SearchDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	results, err := x.ResourceActionsX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, results)
}

type SearchResult struct {
	ID         string `json:"id"`
	ResourceID string `json:"resource_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []*SearchResult, err error) {
	do := x.Db.Model(model.ResourceAction{}).WithContext(ctx)
	if dto.ResourceID != "" {
		do = do.Where(`resource_id = ?`, dto.ResourceID)
	}
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
	}

	results = make([]*SearchResult, 0)
	ctx = common.SetPipe(ctx, common.NewSearchPipe(`id`, `resource_id`, `name`, `code`))
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
