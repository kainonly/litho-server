package index

import (
	"context"
	"fmt"
	"litho-api/common"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type GetUserPhoneDto struct {
	Phone string `query:"phone" vd:"required"`
}

func (x *Controller) GetUserPhoneCode(ctx context.Context, c *app.RequestContext) {
	var dto GetUserPhoneDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if _, err := x.IndexX.GetUserPhoneCode(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) GetUserPhoneCode(ctx context.Context, dto GetUserPhoneDto) (code string, err error) {
	code = help.RandomNumber(6)
	if err = x.TencentX.Sms(ctx, x.V.Sms.Templates["bind_verify"],
		[]string{code},
		[]string{dto.Phone},
	); err != nil {
		return
	}
	key := fmt.Sprintf(`sms-bind:%s`, dto.Phone)
	if exists := x.Captcha.Exists(ctx, key); exists {
		err = common.ErrCodeFrequently
		return
	}
	x.Captcha.Create(ctx, key, code, time.Minute*5)
	return
}
