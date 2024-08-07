package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
	"time"
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

	if _, err := x.IndexX.GetLoginSms(ctx, dto.Phone); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) GetLoginSms(ctx context.Context, phone string) (code string, err error) {
	keyLock := fmt.Sprintf(`phone:%s`, phone)
	if err = x.Locker.Verify(ctx, keyLock, x.V.LoginFailures); err != nil {
		switch err {
		case locker.ErrLockerNotExists:
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
	if n, err = x.Db.Collection("users").
		CountDocuments(ctx, bson.M{"phone": phone, "status": true}); err != nil {
		return
	}
	if n == 0 {
		x.Locker.Update(ctx, keyLock, time.Hour*24)
		err = common.ErrSmsNotExists
		return
	}

	key := fmt.Sprintf(`sms-login:%s`, phone)
	if exists := x.Captcha.Exists(ctx, key); exists {
		err = common.ErrCodeFrequently
		return
	}

	code = help.RandomNumber(6)
	// TODO: Change to values...
	if err = x.TencentX.SmsSend(ctx, x.V.SmsSign, x.V.SmsLoginVerify, []string{code}, []string{phone}); err != nil {
		return
	}

	x.Captcha.Create(ctx, key, code, time.Minute*2)
	return
}
