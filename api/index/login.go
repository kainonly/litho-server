package index

import (
	"context"
	"errors"
	"server/common"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
	"gorm.io/gorm"
)

type LoginDto struct {
	Email    string `json:"email" vd:"required,email"`
	Password string `json:"password" vd:"required,min=8"`
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
	result = new(LoginResult)
	if result, err = x.QueryLoginUser(ctx, func(do *gorm.DB) *gorm.DB {
		return do.Where(`email = ?`, dto.Email)
	}); err != nil {
		err = common.ErrLoginNotExists
		return
	}

	if err = passlib.Verify(dto.Password, result.Password); err != nil {
		if errors.Is(err, passlib.ErrNotMatch) {
			x.Locker.Update(ctx, result.ID, time.Minute*15)
			err = common.ErrLoginInvalid
			return
		}
		return
	}

	if result.AccessToken, err = x.CreateAccessToken(ctx, result.ID); err != nil {
		return
	}
	return
}
