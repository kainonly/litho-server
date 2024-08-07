package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginSmsDto struct {
	Phone string `json:"phone" vd:"required"`
	Code  string `json:"code" vd:"required"`
}

func (x *Controller) LoginSms(ctx context.Context, c *app.RequestContext) {
	var dto LoginSmsDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.IndexX.LoginSms(ctx, dto.Phone, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	data := model.NewLogsetLogin(r.User.ID, string(c.GetHeader(x.V.Ip)), "sms", string(c.UserAgent()))
	if err = x.IndexX.WriteLogsetLogin(context.TODO(), data); err != nil {
		c.Error(err)
		return
	}

	common.SetAccessToken(c, r.AccessToken)
	c.Status(204)
}

func (x *Service) LoginSms(ctx context.Context, phone string, code string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.Logining(ctx, bson.M{"phone": phone, "status": true}); err != nil {
		return
	}

	userId := r.User.ID.Hex()
	key := fmt.Sprintf(`sms-login:%s`, phone)
	if err = x.Captcha.Verify(ctx, key, code); err != nil {
		x.Locker.Update(ctx, userId, x.V.LoginTTL)
		err = common.ErrSmsInvalid
		return
	}

	if r.AccessToken, err = x.CreateAccessToken(ctx, userId); err != nil {
		return
	}
	x.Captcha.Delete(ctx, key)
	return
}
