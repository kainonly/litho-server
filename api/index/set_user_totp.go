package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/totp"
	"github.com/weplanx/server/common"
	"go.mongodb.org/mongo-driver/bson"
)

type SetUserTotp struct {
	Totp string    `json:"totp" vd:"required"`
	Tss  [2]string `json:"tss" vd:"len=2"`
}

func (x *Controller) SetUserTotp(ctx context.Context, c *app.RequestContext) {
	var dto SetUserTotp
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := common.Claims(c)
	if _, err := x.IndexX.SetUserTotp(ctx, claims.UserId, dto.Totp, dto.Tss); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) SetUserTotp(ctx context.Context, userId string, uri string, tss [2]string) (r interface{}, err error) {
	if tss[0] == tss[1] {
		return "", common.ErrTotpInvalid
	}
	var secret string
	if secret, err = x.RDb.Get(ctx, uri).Result(); err != nil {
		return
	}
	otpc := &totp.Totp{
		Secret:  secret,
		Window:  2,
		Counter: 0,
	}
	for _, v := range tss {
		var check bool
		if check, err = otpc.Authenticate(v); err != nil {
			return
		}
		if !check {
			return "", common.ErrTotpInvalid
		}
	}

	if err = x.RDb.Del(ctx, uri).Err(); err != nil {
		return
	}
	return x.SetUser(ctx, userId, bson.M{
		"$set": bson.M{
			"totp": secret,
		},
	})
}
