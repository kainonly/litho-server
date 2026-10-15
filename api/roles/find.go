package roles

import (
	"context"
	"litho-api/common"
	"litho-api/model"
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
	total, data, err := x.RolesX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *bool  `json:"status"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []FindResult, err error) {
	do := x.Db.Model(model.Role{}).WithContext(ctx)
	if dto.Q != "" {
		keyword := dto.GetKeyword()
		do = do.Where(`name like ?`, keyword)
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	results = make([]FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe().
		Select("id", "name", "description", "status"))
	if err = dto.Find(ctx, do.Order(`sort`), &results); err != nil {
		return
	}
	return
}
