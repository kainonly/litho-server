package index

import (
	"context"
	"server/common"
	"server/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
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
	ID             string     `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	LoginAt        *time.Time `json:"login_at"`
	DepartmentID   string     `json:"-"`
	Department     string     `json:"department"`
	DepartmentType *int16     `json:"department_type"`
	RoleID         string     `json:"-"`
	Role           string     `json:"role"`
	Cabs           []string   `json:"cabs"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	Name           string     `json:"name"`
	Avatar         string     `json:"avatar"`
	Sessions       int32      `json:"sessions"`
}

func (x *Service) GetUser(ctx context.Context, userId string) (result *UserResult, err error) {
	var data *model.User
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, userId).
		Take(&data).Error; err != nil {
		return
	}

	result = &UserResult{
		ID:           data.ID,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
		LoginAt:      data.LoginAt,
		DepartmentID: data.DepartmentID,
		RoleID:       data.RoleID,
		Email:        data.Email,
		Phone:        data.Phone,
		Name:         data.Name,
		Avatar:       data.Avatar,
		Sessions:     data.Sessions,
	}

	result.Department = "SYS"
	if result.DepartmentID != "0" {
		var dept *model.Department
		if err = x.Db.Model(model.Department{}).WithContext(ctx).
			Select([]string{`id`, `name`, `type`}).
			Where(`id = ?`, result.DepartmentID).
			Take(&dept).Error; err != nil {
			return
		}
		result.Department = dept.Name
		result.DepartmentType = dept.Type
	}
	result.Role = "无"
	if result.RoleID != "0" {
		var role *model.Role
		if err = x.Db.Model(model.Role{}).WithContext(ctx).
			Select([]string{`id`, `name`, `strategy`}).
			Where(`id = ?`, result.RoleID).
			Take(&role).Error; err != nil {
			return
		}
		result.Role = role.Name
		result.Cabs = role.Strategy.Permissions
	}
	return
}
