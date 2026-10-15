package users

import (
	"context"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
)

type UpdateDto struct {
	ID       string  `json:"id" vd:"required"`
	Phone    string  `json:"phone" vd:"required"`
	Name     string  `json:"name" vd:"required"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
	Avatar   *string `json:"avatar"`
	Status   *bool   `json:"status"`
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

	common.SetTracking(c, Key, "update", dto.ID)
	c.JSON(200, help.Ok())
}

func (x *Service) Update(ctx context.Context, user *common.IAMUser, dto UpdateDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}

	data := M{
		"update_time": time.Now(),
		"name":        dto.Name,
		"phone":       dto.Phone,
		"email":       dto.Email,
	}

	if dto.Password != "" {
		data["password"], _ = passlib.Hash(dto.Password)
	}
	if dto.Status != nil {
		data["status"] = *dto.Status
	}
	if dto.Avatar != nil {
		data["avatar"] = *dto.Avatar
	}

	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, dto.ID).
		Updates(data).Error; err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
