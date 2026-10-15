package index

import (
	"context"
	"fmt"
	"litho-api/common"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"gorm.io/gorm"
)

type LoginSmsDto struct {
	Phone string `json:"phone" vd:"required"`
	Code  string `json:"code" vd:"required,len=6"`
}

func (x *Controller) LoginSms(ctx context.Context, c *app.RequestContext) {
	var dto LoginSmsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user, err := x.IndexX.LoginSms(ctx, dto)
	if err != nil {
		c.Error(err)
		return
	}

	x.IndexX.SetAccessToken(c, user.AccessToken)
	x.LoginAudit(c, user, 2)
	c.JSON(200, help.Ok())
}

func (x *Service) LoginSms(ctx context.Context, dto LoginSmsDto) (result *LoginResult, err error) {
	result = new(LoginResult)
	if result, err = x.QueryLogin(ctx, func(do *gorm.DB) *gorm.DB {
		return do.Where("phone = ?", dto.Phone)
	}); err != nil {
		err = common.ErrLoginNotExists
		return
	}

	key := fmt.Sprintf(`sms-login:%s`, dto.Phone)
	if err = x.Captcha.Verify(ctx, key, dto.Code); err != nil {
		if _, err = x.Locker.Increment(ctx, result.ID, time.Minute*2); err != nil {
			return
		}
		err = common.ErrSmsInvalid
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
