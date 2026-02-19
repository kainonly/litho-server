package users

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
)

type CreateDto struct {
	ID       string `json:"-"`
	Active   *bool  `json:"active" vd:"required"`
	Email    string `json:"email" vd:"required,email"`
	Password string `json:"password" vd:"required,min=8"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	dto.ID = help.SID()
	if err := x.UsersX.Create(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.User{
		ID:     dto.ID,
		Active: dto.Active,
		Email:  dto.Email,
	}
	if dto.Phone != "" {
		data.Phone = dto.Phone
	}
	if dto.Name != "" {
		data.Name = dto.Name
	}
	if dto.Avatar != "" {
		data.Avatar = dto.Avatar
	}
	data.Password, _ = passlib.Hash(dto.Password)
	if err = x.Db.WithContext(ctx).Create(data).Error; err != nil {
		return
	}
	if err = x.Db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
