package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
)

type SetUserPhone struct {
	Phone string `json:"phone" vd:"required"`
	Code  string `json:"code" vd:"required"`
}

func (x *Controller) SetUserPhone(ctx context.Context, c *app.RequestContext) {
	var dto SetUserPhone
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := common.Claims(c)
	if _, err := x.IndexX.SetUserPhone(ctx, claims.UserId, dto.Phone, dto.Code); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) SetUserPhone(ctx context.Context, userId string, phone string, code string) (r interface{}, err error) {
	key := fmt.Sprintf(`sms-bind:%s`, phone)
	if err = x.Captcha.Verify(ctx, key, code); err != nil {
		err = common.ErrSmsInvalid
		return
	}

	x.Captcha.Delete(ctx, key)
	return x.SetUser(ctx, userId, bson.M{
		"$set": bson.M{
			"phone": phone,
		},
	})
}
