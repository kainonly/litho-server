package users

import (
	"context"
	"strconv"

	"server/common"
	"server/model"

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
	total, data, err := x.UsersX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(&model.User{}).WithContext(ctx)
	if dto.Q != "" {
		do = do.Where(`name like ? OR email like ? OR phone like ?`, dto.GetKeyword(), dto.GetKeyword(), dto.GetKeyword())
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	results = make([]*FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
