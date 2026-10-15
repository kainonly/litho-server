package roles

import (
	"context"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
)

func (x *Controller) GetSuper(ctx context.Context, c *app.RequestContext) {
	data, err := x.RolesX.GetSuper(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type GetSuperResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (x *Service) GetSuper(ctx context.Context) (result GetSuperResult, err error) {
	if err = x.Db.Model(model.Role{}).WithContext(ctx).
		Select("id", "name").
		Where("name = ?", `超级管理员`).
		Take(&result).Error; err != nil {
		return
	}
	return
}
