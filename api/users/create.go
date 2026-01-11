package users

import (
	"context"

	"server/common"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
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
		c.JSON(400, utils.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(*common.IAMUser)
	if err := x.UsersX.Create(ctx, user, dto); err != nil {
		c.JSON(500, utils.H{"error": err.Error()})
		return
	}
	c.JSON(201, utils.H{"ok": true})
}

func (x *Service) Create(ctx context.Context, user *common.IAMUser, dto CreateDto) (err error) {
	return
}
