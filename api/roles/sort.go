package roles

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type SortDto struct {
	IDs []string `json:"ids"`
}

func (x *Controller) Sort(ctx context.Context, c *app.RequestContext) {
	var dto SortDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RolesX.Sort(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "sort", "0")
	c.JSON(200, help.Ok())
}

func (x *Service) Sort(ctx context.Context, user *common.IAMUser, dto SortDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	return x.Db.Transaction(func(tx *gorm.DB) (errX error) {
		for i, id := range dto.IDs {
			updates := M{
				"update_time": time.Now(),
				"sort":        i,
			}

			if errX = tx.Model(model.Role{}).WithContext(ctx).
				Where(`id = ?`, id).
				Updates(updates).Error; errX != nil {
				return
			}
		}
		return
	})
}
