package users

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type CreateDto struct {
	ID       string `json:"-"`
	Active   *bool  `json:"active" vd:"required"`
	Email    string `json:"email" vd:"required,email"`
	Phone    string `json:"phone" vd:"required"`
	Name     string `json:"name" vd:"required"`
	Password string `json:"password" vd:"required,min=6"`
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
		ID:       dto.ID,
		Active:   *dto.Active,
		Email:    dto.Email,
		Phone:    dto.Phone,
		Name:     dto.Name,
		Password: dto.Password,
		Avatar:   dto.Avatar,
	}
	if err = x.Db.WithContext(ctx).
		Create(&data).Error; err != nil {
		return
	}
	return
}
