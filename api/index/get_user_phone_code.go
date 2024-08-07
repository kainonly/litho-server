package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"github.com/weplanx/server/common"
	"time"
)

type GetUserPhone struct {
	Phone string `query:"phone" vd:"required"`
}

func (x *Controller) GetUserPhoneCode(ctx context.Context, c *app.RequestContext) {
	var dto GetUserPhone
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if _, err := x.IndexX.GetUserPhoneCode(ctx, dto.Phone); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) GetUserPhoneCode(ctx context.Context, phone string) (code string, err error) {
	code = help.RandomNumber(6)
	// TODO: Change to values
	if err = x.TencentX.SmsSend(ctx, x.V.SmsSign, x.V.SmsPhoneBind, []string{code}, []string{phone}); err != nil {
		return
	}
	key := fmt.Sprintf(`sms-bind:%s`, phone)
	if exists := x.Captcha.Exists(ctx, key); exists {
		err = common.ErrCodeFrequently
		return
	}
	x.Captcha.Create(ctx, key, code, time.Minute*2)
	return
}
