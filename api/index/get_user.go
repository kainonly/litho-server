package index

import (
	"context"
	"errors"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	user := common.GetIAM(c)
	data, err := x.IndexX.GetUser(ctx, user.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

type UserResult struct {
	ID          string    `json:"id"`
	RoleIDs     []string  `json:"role_ids"`
	Roles       []string  `json:"roles"`
	Permissions []any     `json:"permissions"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
	Sessions    int32     `json:"sessions"`
	CreateTime  time.Time `json:"create_time"`
	UpdateTime  time.Time `json:"update_time"`
}

func (x *Service) GetUser(ctx context.Context, userId string) (result *UserResult, err error) {
	var data *model.User
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, userId).
		Take(&data).Error; err != nil {
		return
	}

	result = &UserResult{
		ID:         data.ID,
		RoleIDs:    make([]string, 0),
		Roles:      make([]string, 0),
		Email:      data.Email,
		Phone:      data.Phone,
		Name:       data.Name,
		Avatar:     data.Avatar,
		Sessions:   data.Sessions,
		CreateTime: data.CreateTime,
		UpdateTime: data.UpdateTime,
	}

	// 用户只关联一个角色；响应仍保留数组形式以兼容前端。
	var userRole model.UserRole
	if err = x.Db.Model(model.UserRole{}).WithContext(ctx).
		Select(`role_id`).
		Where(`user_id = ?`, userId).
		Take(&userRole).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
		return
	} else if err != nil {
		return
	}
	result.RoleIDs = append(result.RoleIDs, userRole.RoleID)

	var role model.Role
	if err = x.Db.Model(model.Role{}).WithContext(ctx).
		Select([]string{`id`, `name`, `strategy`}).
		Where(`id = ?`, userRole.RoleID).
		Take(&role).Error; err != nil {
		return
	}
	result.Roles = append(result.Roles, role.Name)
	var b []byte
	if b, err = sonic.Marshal(role.Strategy); err != nil {
		return
	}
	var strategy common.IAMStrategy
	if err = sonic.Unmarshal(b, &strategy); err != nil {
		return
	}
	result.Permissions = strategy.Permissions
	return
}
