package index

import (
	"context"
	"fmt"
	"litho-api/common"
	"litho-api/model"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
	"github.com/kainonly/go/locker"
)

type GetLoginSmsDto struct {
	Phone string `query:"phone" vd:"required"`
}

func (x *Controller) GetLoginSms(ctx context.Context, c *app.RequestContext) {
	var dto GetLoginSmsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if _, err := x.IndexX.GetLoginSms(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) GetLoginSms(ctx context.Context, dto GetLoginSmsDto) (code string, err error) {
	keyLock := fmt.Sprintf(`phone:%s`, dto.Phone)
	if err = x.Locker.Check(ctx, keyLock, 5); err != nil {
		switch err {
		case locker.ErrNotExists:
			err = nil
			break
		case locker.ErrLocked:
			err = common.ErrLoginMaxFailures
			return
		default:
			return
		}
	}

	var n int64
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`phone = ?`, dto.Phone).
		Where(`status = ?`, true).
		Count(&n).Error; err != nil {
		return
	}

	if n == 0 {
		if _, err = x.Locker.Increment(ctx, keyLock, time.Hour*24); err != nil {
			return
		}
		err = common.ErrSmsNotExists
		return
	}

	key := fmt.Sprintf(`sms-login:%s`, dto.Phone)
	if exists := x.Captcha.Exists(ctx, key); exists {
		err = common.ErrCodeFrequently
		return
	}

	code = help.RandomNumber(6)
	if err = x.TencentX.Sms(ctx, x.V.Sms.Templates["login_verify"],
		[]string{code, "2"},
		[]string{dto.Phone},
	); err != nil {
		return
	}

	x.Captcha.Create(ctx, key, code, time.Minute*2)
	return
}
