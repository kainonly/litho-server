package products

import (
	"context"
	"server/common"
	"server/model"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type FindDto struct {
	common.FindDto
	OrgID string `query:"org_id,omitempty"`
}

func (x *Controller) Find(ctx context.Context, c *app.RequestContext) {
	var dto FindDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	total, data, err := x.ProductsX.Find(ctx, user, dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("x-total", strconv.Itoa(int(total)))
	c.JSON(200, data)
}

type FindResult struct {
	ID          string  `json:"id"`
	OrgID       string  `json:"org_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	Active      bool    `json:"active"`
	Thumbnail   string  `json:"thumbnail"`
}

func (x *Service) Find(ctx context.Context, user *common.IAMUser, dto FindDto) (total int64, results []*FindResult, err error) {
	do := x.Db.Model(&model.Product{}).WithContext(ctx)
	if dto.OrgID != "" {
		do = do.Where(`org_id = ?`, dto.OrgID)
	}
	if dto.Q != "" {
		do = do.Where(`name like ?`, dto.GetKeyword())
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
