package index

import (
	"context"
	"errors"
	"litho-api/common"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/passlib"
	"gorm.io/gorm"
)

type LoginDto struct {
	Username string `json:"username" vd:"required"`
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
	x.LoginAudit(c, user, 1)
	c.JSON(200, help.Ok())
}

func (x *Service) Login(ctx context.Context, dto LoginDto) (result *LoginResult, err error) {
	result = new(LoginResult)
	if result, err = x.QueryLogin(ctx, func(do *gorm.DB) *gorm.DB {
		return do.Where("phone = ?", dto.Username)
	}); err != nil {
		err = common.ErrLoginNotExists
		return
	}

	if result.Password == "" {
		err = help.E(0, `登录账号需要初始化，请使用短信登录并重置个人密码`)
		return
	}

	if err = passlib.Verify(dto.Password, result.Password); err != nil {
		if errors.Is(err, passlib.ErrNotMatch) {
			if _, err = x.Locker.Increment(ctx, result.ID, time.Minute*15); err != nil {
				return
			}
			err = common.ErrLoginInvalid
			return
		}
		return
	}

	if err = x.UpdateLogin(ctx, result.ID); err != nil {
		return
	}
	if result.AccessToken, err = x.CreateAccessToken(ctx, result.ID); err != nil {
		return
	}
	return
}
