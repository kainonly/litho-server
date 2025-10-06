package users

import (
	"context"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passlib"
)

type CreateDto struct {
	ID       string `json:"-"`
	Email    string `json:"email"`
	Password string `json:"password" vd:"required"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	dto.ID = help.SID()
	user := common.GetIAM(c)
	if err := x.UsersX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := &model.User{
		ID:    dto.ID,
		Email: dto.Email,
	}
	data.Password, _ = passlib.Hash(dto.Password)
	if err = x.Db.WithContext(ctx).Create(data).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
