package users

import (
	"context"
	"litho-api/common"
	"litho-api/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type SetRolesDto struct {
	IDs   []string `json:"ids" vd:"required"`
	Roles []string `json:"roles" vd:"required"`
}

func (x *Controller) SetRoles(ctx context.Context, c *app.RequestContext) {
	var dto SetRolesDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if err := x.UsersX.SetRoles(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	common.SetTracking(c, Key, "set_roles", dto.IDs...)
	c.JSON(200, help.Ok())
}

func (x *Service) SetRoles(ctx context.Context, user *common.IAMUser, dto SetRolesDto) (err error) {
	if err = x.CheckAccess(ctx, user); err != nil {
		return
	}
	if len(dto.Roles) != 1 {
		return help.E(0, `操作失败，每个用户只能设置一个权限组`)
	}

	if err = x.Db.Transaction(func(tx *gorm.DB) (errX error) {
		if errX = tx.WithContext(ctx).
			Where(`user_id in (?)`, dto.IDs).
			Delete(model.UserRole{}).Error; errX != nil {
			return
		}
		records := make([]model.UserRole, 0)
		for _, userId := range dto.IDs {
			records = append(records, model.UserRole{
				UserID: userId,
				RoleID: dto.Roles[0],
			})
		}
		if len(records) != 0 {
			return tx.WithContext(ctx).Create(records).Error
		}
		return
	}); err != nil {
		return
	}
	return x.RefreshCache(ctx)
}
