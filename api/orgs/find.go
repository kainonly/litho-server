package orgs

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto

	Type int16 `query:"type"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.OrgsX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID        string `json:"id"`
	Status    bool   `json:"status"`
	Pid       string `json:"pid"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Type      int16  `json:"type"`
	Company   bool   `json:"company"`
	ContactID string `json:"contact_id"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(&model.Org{}).WithContext(ctx)

	// 查询条件过滤
	if dto.Type != 0 {
		do = do.Where("type = ?", dto.Type)
	}
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
	}

	if err = do.Count(&total).Error; err != nil {
		return
	}

	if len(dto.Sort) == 0 {
		do = do.Order(`type`).Order("create_time desc")
	}

	results = make([]*FindResult, 0)
	ctx = common.SetPipe(ctx, common.NewFindPipe().SkipTs())
	if err = dto.Find(ctx, do, &results); err != nil {
		return
	}
	return
}
