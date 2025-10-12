package schedulers

import (
	"fmt"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"golang.org/x/net/context"
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
	data, err := x.SchedulesX.Search(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type SearchResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (x *Service) Search(ctx context.Context, user *common.IAMUser, dto SearchDto) (results []SearchResult, err error) {
	do := x.Db.Model(model.Scheduler{}).WithContext(ctx)

	if dto.Q != "" {
		do = do.Where(`name like ?`, fmt.Sprintf(`%%%s%%`, dto.Q))
	}

	results = make([]SearchResult, 0, len(results))
	ctx = common.SetPipe(ctx, common.NewSearchPipe(`id`, `name`).SkipAsync())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
