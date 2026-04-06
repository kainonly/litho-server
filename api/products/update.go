package products

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type UpdateDto struct {
	ID          string  `json:"id" vd:"required"`
	Name        string  `json:"name" vd:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" vd:"required,gt=0"`
	Stock       int32   `json:"stock" vd:"required,gte=0"`
	Active      *bool   `json:"active" vd:"required"`
	Thumbnail   string  `json:"thumbnail"`
}

const IUpdate = "更新"

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.ProductsX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	return x.Db.Model(model.Product{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(common.M{
			`updated_at`:  time.Now(),
			`name`:        dto.Name,
			`description`: dto.Description,
			`price`:       dto.Price,
			`stock`:       dto.Stock,
			`active`:      *dto.Active,
			`thumbnail`:   dto.Thumbnail,
		}).Error
}
