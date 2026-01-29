package resource_actions

import (
	"context"
	"strconv"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
	ResourceID string `json:"resource_id" query:"resource_id"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.ResourceActionsX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID         string `json:"id"`
	ResourceID string `json:"resource_id"`
	Sort       int16  `json:"sort"`
	Active     bool   `json:"active"`
	Name       string `json:"name"`
	Code       string `json:"code"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(model.ResourceAction{}).WithContext(ctx)
	if dto.ResourceID != "" {
		do = do.Where(`resource_id = ?`, dto.ResourceID)
	}
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	results = make([]*FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe())
	if err = dto.Find(ctx, do.Order(`sort`), &results); err != nil {
		return
	}
	return
}
