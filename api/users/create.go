package users

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/snowflake"
)

type CreateDto struct {
	Active   *bool  `json:"active,omitempty"`
	Email    string `json:"email" vd:"required,email"`
	Phone    string `json:"phone" vd:"required"`
	Name     string `json:"name" vd:"required"`
	Password string `json:"password" vd:"required,min=6"`
	Avatar   string `json:"avatar,omitempty"`
}

func (x *Controller) Create(ctx context.Context, c *app.RequestContext) {
	var dto CreateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.UsersX.Create(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(201, map[string]any{"ok": true})
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	data := model.User{
		ID:       snowflake.GenerateString(),
		Email:    dto.Email,
		Phone:    dto.Phone,
		Name:     dto.Name,
		Password: dto.Password, // 实际应用需要加密
		Avatar:   dto.Avatar,
	}
	if dto.Active != nil {
		data.Active = *dto.Active
	} else {
		data.Active = true
	}
	return x.Db.Create(&data).Error
}
