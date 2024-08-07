package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/weplanx/go/passlib"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SetUserPassword struct {
	Old      string `json:"old" vd:"min=8"`
	Password string `json:"password" vd:"min=8"`
}

func (x *Controller) SetUserPassword(ctx context.Context, c *app.RequestContext) {
	var dto SetUserPassword
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	claims := common.Claims(c)
	if _, err := x.IndexX.SetUserPassword(ctx, claims.UserId, dto.Old, dto.Password); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func (x *Service) SetUserPassword(ctx context.Context, userId string, old string, password string) (r interface{}, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return
	}

	if err = passlib.Verify(old, user.Password); err != nil {
		if err == passlib.ErrNotMatch {
			err = errors.NewPublic(passlib.ErrNotMatch.Error())
			return
		}
	}

	var hash string
	if hash, err = passlib.Hash(password); err != nil {
		return
	}
	return x.SetUser(ctx, userId, bson.M{
		"$set": bson.M{
			"password": hash,
		},
	})
}
