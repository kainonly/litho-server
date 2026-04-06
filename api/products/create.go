package products

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID          string  `json:"-"`
	OrgID       string  `json:"org_id" vd:"required"`
	Name        string  `json:"name" vd:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" vd:"required,gt=0"`
	Stock       int32   `json:"stock" vd:"required,gte=0"`
	Active      *bool   `json:"active" vd:"required"`
	Thumbnail   string  `json:"thumbnail"`
}

const ICreate = "新增"

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.ProductsX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.Product{
		ID:          dto.ID,
		OrgID:       dto.OrgID,
		Name:        dto.Name,
		Description: dto.Description,
		Price:       dto.Price,
		Stock:       dto.Stock,
		Active:      dto.Active,
		Thumbnail:   dto.Thumbnail,
	}
	return x.Db.WithContext(ctx).Create(&data).Error
}
