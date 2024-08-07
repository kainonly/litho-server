package index

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/url"
	"time"
)

func (x *Controller) GetUserTotp(ctx context.Context, c *app.RequestContext) {
	claims := common.Claims(c)
	uri, err := x.IndexX.GetUserTotp(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, M{
		"totp": uri,
	})
}

func (x *Service) GetUserTotp(ctx context.Context, userId string) (uri string, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return
	}
	random := make([]byte, 10)
	if _, err = rand.Read(random); err != nil {
		return
	}
	secret := base32.StdEncoding.EncodeToString(random)
	var u *url.URL
	if u, err = url.Parse("otpauth://totp"); err != nil {
		return
	}
	u.Path += "/" + url.PathEscape(x.V.Namespace) + ":" + url.PathEscape(user.Email)
	params := url.Values{}
	params.Add("secret", secret)
	params.Add("issuer", x.V.Namespace)
	u.RawQuery = params.Encode()
	uri = u.String()

	if err = x.RDb.Set(ctx, uri, secret, time.Minute*5).Err(); err != nil {
		return
	}
	return
}
