package index

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/passlib"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
)

type ForgetResetDto struct {
	Email    string `json:"email" vd:"email"`
	Code     string `json:"code" vd:"required"`
	Password string `json:"password" vd:"required"`
}

func (x *Controller) ForgetReset(ctx context.Context, c *app.RequestContext) {
	var dto ForgetResetDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.IndexX.ForgetReset(ctx, dto); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) ForgetReset(ctx context.Context, dto ForgetResetDto) (err error) {
	key := fmt.Sprintf(`forget:%s`, dto.Email)
	if err = x.Captcha.Verify(ctx, key, dto.Code); err != nil {
		err = common.ErrEmailInvalid
		return
	}

	filter := bson.M{"email": dto.Email, "status": true}
	password, _ := passlib.Hash(dto.Password)
	data := bson.M{"$set": bson.M{"password": password}}
	if _, err = x.Db.Collection("users").UpdateOne(ctx, filter, data); err != nil {
		return
	}

	x.Captcha.Delete(ctx, key)
	return
}
