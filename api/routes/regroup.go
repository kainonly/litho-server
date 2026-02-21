package routes

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type RegroupDto struct {
	Update RegroupUpdate `json:"update" vd:"required"`
	Sorts  [][]string    `json:"sorts" vd:"required"`
}

type RegroupUpdate struct {
	Changed *bool  `json:"changed" vd:"required"`
	ID      string `json:"id" vd:"required"`
	Pid     string `json:"pid" vd:"required_if=Changed true"`
}

const IRegroup = "重新分组"

func (x *Controller) Regroup(ctx context.Context, c *app.RequestContext) {
	var dto RegroupDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.RoutesX.Regroup(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Regroup(ctx context.Context, user *common.IAMUser, dto RegroupDto) (err error) {
	return x.Db.Transaction(func(tx *gorm.DB) (errX error) {
		if *dto.Update.Changed {
			if errX = tx.Model(model.Route{}).WithContext(ctx).
				Where(`id = ?`, dto.Update.ID).
				Update("pid", dto.Update.Pid).
				Error; errX != nil {
				return
			}
		}

		for _, sort := range dto.Sorts {
			for index, id := range sort {
				if errX = tx.Model(model.Route{}).WithContext(ctx).
					Where(`id = ?`, id).
					Update("sort", index).
					Error; errX != nil {
					return
				}
			}
		}

		return
	})
}
