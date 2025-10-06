package users

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passlib"
)

type UpdateDto struct {
	ID       string `json:"id" vd:"required"`
	Password string `json:"password"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.UsersX.Update(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	data := M{
		"update_time": time.Now(),
	}

	if dto.Password != "" {
		data["password"], _ = passlib.Hash(dto.Password)
	}

	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(data).
		Error; err != nil {
		return
	}

	return x.RefreshCache(ctx)
}
