package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/totp"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginTotpDto struct {
	Email string `json:"email" vd:"email"`
	Code  string `json:"code" vd:"required"`
}

func (x *Controller) LoginTotp(ctx context.Context, c *app.RequestContext) {
	var dto LoginTotpDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.IndexX.LoginTotp(ctx, dto.Email, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	data := model.NewLogsetLogin(r.User.ID, string(c.GetHeader(x.V.Ip)), "totp", string(c.UserAgent()))
	if err = x.IndexX.WriteLogsetLogin(ctx, data); err != nil {
		c.Error(err)
		return
	}

	common.SetAccessToken(c, r.AccessToken)
	c.Status(204)
}

func (x *Service) LoginTotp(ctx context.Context, email string, code string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.Logining(ctx, bson.M{"email": email, "status": true}); err != nil {
		return
	}

	userId := r.User.ID.Hex()
	otpc := &totp.Totp{
		Secret:  r.User.Totp,
		Window:  1,
		Counter: 0,
	}
	var check bool
	if check, err = otpc.Authenticate(code); err != nil {
		return
	}
	if !check {
		x.Locker.Update(ctx, userId, x.V.LoginTTL)
		err = common.ErrLoginInvalid
		return
	}

	if r.AccessToken, err = x.CreateAccessToken(ctx, userId); err != nil {
		return
	}
	return
}
