package orders

import (
	"context"
	"fmt"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type CreateItemDto struct {
	ProductID   string  `json:"product_id" vd:"required"`
	ProductName string  `json:"product_name" vd:"required"`
	Price       float64 `json:"price" vd:"required,gt=0"`
	Quantity    int32   `json:"quantity" vd:"required,gt=0"`
	Subtotal    float64 `json:"subtotal" vd:"required,gt=0"`
}

type CreateDto struct {
	ID           string          `json:"-"`
	DepartmentID string          `json:"department_id" vd:"required"`
	UserID       string          `json:"user_id" vd:"required"`
	ScheduledAt  time.Time       `json:"scheduled_at" vd:"required"`
	Remark       string          `json:"remark"`
	Items        []CreateItemDto `json:"items" vd:"required,gt=0,dive"`
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
	if err := x.OrdersX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) error {
	var amount float64
	for _, item := range dto.Items {
		amount += item.Subtotal
	}

	return x.Db.Transaction(func(tx *gorm.DB) (errX error) {
		order := model.Order{
			ID:           dto.ID,
			DepartmentID: dto.DepartmentID,
			UserID:       dto.UserID,
			No:           fmt.Sprintf("ORD%s", dto.ID),
			Amount:       amount,
			Status:       0,
			ScheduledAt:  dto.ScheduledAt,
			Remark:       dto.Remark,
		}
		if errX = tx.WithContext(ctx).Create(&order).Error; errX != nil {
			return
		}

		items := make([]model.OrderItem, 0, len(dto.Items))
		for _, v := range dto.Items {
			qty := v.Quantity
			items = append(items, model.OrderItem{
				ID:          help.SID(),
				OrderID:     dto.ID,
				ProductID:   v.ProductID,
				ProductName: v.ProductName,
				Price:       v.Price,
				Quantity:    &qty,
				Subtotal:    v.Subtotal,
			})
		}
		if errX = tx.WithContext(ctx).Create(&items).Error; errX != nil {
			return
		}
		return
	})
}
