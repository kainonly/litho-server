package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/passlib"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
)

type LoginDto struct {
	Email    string `json:"email" vd:"email"`
	Password string `json:"password" vd:"min=8"`
}

func (x *Controller) Login(ctx context.Context, c *app.RequestContext) {
	var dto LoginDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	r, err := x.IndexX.Login(ctx, dto.Email, dto.Password)
	if err != nil {
		c.Error(err)
		return
	}

	data := model.NewLogsetLogin(r.User.ID, string(c.GetHeader(x.V.Ip)), "email", string(c.UserAgent()))
	if err = x.IndexX.WriteLogsetLogin(ctx, data); err != nil {
		c.Error(err)
		return
	}

	common.SetAccessToken(c, r.AccessToken)
	c.Status(204)
}

func (x *Service) Login(ctx context.Context, username string, password string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.Logining(ctx, bson.M{"email": username, "status": true}); err != nil {
		return
	}

	userId := r.User.ID.Hex()
	if err = passlib.Verify(password, r.User.Password); err != nil {
		if err == passlib.ErrNotMatch {
			x.Locker.Update(ctx, userId, x.V.LoginTTL)
			err = common.ErrLoginInvalid
			return
		}
		return
	}

	if r.AccessToken, err = x.CreateAccessToken(ctx, userId); err != nil {
		return
	}
	return
}
