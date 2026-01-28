package users

import (
	"context"
	"time"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
)

type UpdateDto struct {
	ID       string `json:"id" vd:"required"`
	Active   *bool  `json:"active" vd:"required"`
	Email    string `json:"email" vd:"required,email"`
	Phone    string `json:"phone" vd:"required"`
	Name     string `json:"name" vd:"required"`
	Password string `json:"password" vd:"omitempty,min=6"`
	Avatar   string `json:"avatar"`
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
	updates := common.M{
		`updated_at`: time.Now(),
		`active`:     *dto.Active,
		`email`:      dto.Email,
		`phone`:      dto.Phone,
		`name`:       dto.Name,
		`avatar`:     dto.Avatar,
	}
	if dto.Password != "" {
		updates[`password`], _ = passlib.Hash(dto.Password)
	}
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(updates).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
