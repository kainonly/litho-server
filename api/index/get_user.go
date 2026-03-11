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
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LoginAt   *time.Time `json:"login_at"`
	OrgID     string     `json:"-"`
	Org       string     `json:"org"`
	OrgType   *int16     `json:"org_type"`
	RoleID    string     `json:"-"`
	Role      string     `json:"role"`
	Cabs      []string   `json:"cabs"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Name      string     `json:"name"`
	Avatar    string     `json:"avatar"`
	Sessions  int32      `json:"sessions"`
}

func (x *Service) GetUser(ctx context.Context, userId string) (result *UserResult, err error) {
	var data *model.User
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, userId).
		Take(&data).Error; err != nil {
		return
	}

	result = &UserResult{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		LoginAt:   data.LoginAt,
		OrgID:     data.OrgID,
		RoleID:    data.RoleID,
		Email:     data.Email,
		Phone:     data.Phone,
		Name:      data.Name,
		Avatar:    data.Avatar,
		Sessions:  data.Sessions,
	}

	result.Org = "SYS"
	if result.OrgID != "0" {
		var org *model.Org
		if err = x.Db.Model(model.Org{}).WithContext(ctx).
			Select([]string{`id`, `name`, `type`}).
			Where(`id = ?`, result.OrgID).
			Take(&org).Error; err != nil {
			return
		}
		result.Org = org.Name
		result.OrgType = org.Type
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
		result.Cabs = role.Strategy.Caps
	}
	return
}
