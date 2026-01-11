package users

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type SetRolesDto struct {
	IDs    []string `json:"ids" vd:"required,min=1"`
	RoleID string   `json:"role_id" vd:"required"`
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

	c.JSON(200, help.Ok())
}

func (x *Service) SetRoles(ctx context.Context, user *common.IAMUser, dto SetRolesDto) (err error) {
	return
}
