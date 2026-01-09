package users

import (
	"context"

	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
)

type UpdateDto struct {
	ID       string `json:"id" vd:"required"`
	Active   *bool  `json:"active,omitempty"`
	Email    string `json:"email,omitempty" vd:"omitempty,email"`
	Phone    string `json:"phone,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty" vd:"omitempty,min=6"`
	Avatar   string `json:"avatar,omitempty"`
}

func (x *Controller) Update(ctx context.Context, c *app.RequestContext) {
	var dto UpdateDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.UsersX.Update(ctx, user, dto); err != nil {
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}
	c.JSON(200, map[string]any{"ok": true})
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	updates := make(map[string]interface{})
	if dto.Active != nil {
		updates["active"] = *dto.Active
	}
	if dto.Email != "" {
		updates["email"] = dto.Email
	}
	if dto.Phone != "" {
		updates["phone"] = dto.Phone
	}
	if dto.Name != "" {
		updates["name"] = dto.Name
	}
	if dto.Password != "" {
		updates["password"] = dto.Password // 实际应用需要加密
	}
	if dto.Avatar != "" {
		updates["avatar"] = dto.Avatar
	}
	return x.Db.Model(&model.User{}).Where("id = ?", dto.ID).Updates(updates).Error
}
