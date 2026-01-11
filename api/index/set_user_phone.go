package index

import (
	"context"
	"fmt"
	"server/common"
	"server/model"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/kainonly/go/help"
)

type SetUserPhoneDto struct {
	Phone string `json:"phone" vd:"required"`
	Code  string `json:"code" vd:"required"`
}

func (x *Controller) SetUserPhone(ctx context.Context, c *app.RequestContext) {
	var dto SetUserPhoneDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	user := common.GetIAM(c)
	if _, err := x.IndexX.SetUserPhone(ctx, user, dto); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, help.Ok())
}

func (x *Service) SetUserPhone(ctx context.Context, user *common.IAMUser, dto SetUserPhoneDto) (r interface{}, err error) {
	key := fmt.Sprintf(`sms-bind:%s`, dto.Phone)
	if err = x.Captcha.Verify(ctx, key, dto.Code); err != nil {
		err = common.ErrSmsInvalid
		return
	}

	x.Captcha.Delete(ctx, key)
	if err = x.Db.Model(model.User{}).WithContext(ctx).
		Where(`id = ?`, user.ID).
		Update("phone", dto.Phone).Error; err != nil {
		return
	}

	return
}
