package index

import (
	"context"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
	"gorm.io/gorm"
)

type LoginDto struct {
	Email    string `json:"email" vd:"email"`
	Password string `json:"password" vd:"min=8"`
}

func (x *Controller) Login(ctx context.Context, c *app.RequestContext) {
	var dto LoginDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user, err := x.IndexX.Login(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	x.IndexX.SetAccessToken(c, user.AccessToken)
	c.JSON(200, help.Ok())
}

func (x *Service) Login(ctx context.Context, dto LoginDto) (result *LoginResult, err error) {
	var user *model.User
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`email = ?`, dto.Email).
		Take(&user).Error; err != nil {
		return
	}

	if user == nil {
		err = help.E(0, ``)
		return
	}

	if !*user.Status {
		err = help.E(0, ``)
		return
	}

	if err = passlib.Verify(dto.Password, user.Password); err != nil {
		return
	}

	result = &LoginResult{User: user}

	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, user.ID).
		Update(`sessions`, gorm.Expr(`sessions + ?`, 1)).
		Error; err != nil {
		return
	}

	if result.AccessToken, err = x.CreateAccessToken(ctx, result.ID); err != nil {
		return
	}
	return
}
