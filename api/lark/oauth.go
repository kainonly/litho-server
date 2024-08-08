package lark

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthDto struct {
	Code  string   `query:"code" vd:"required"`
	State StateDto `query:"state"`
}

type StateDto struct {
	Action string `json:"action,omitempty"`
	Locale string `json:"locale,omitempty"`
}

func (x *Controller) OAuth(ctx context.Context, c *app.RequestContext) {
	var dto OAuthDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}
	userData, err := x.LarkX.GetUserAccessToken(ctx, dto.Code)
	if err != nil {
		c.Error(err)
		return
	}

	switch dto.State.Action {
	case "link":
		ts := c.Cookie("TOKEN")
		if ts == nil {
			c.JSON(401, utils.H{
				"code":    0,
				"message": common.ErrAuthenticationExpired.Error(),
			})
			return
		}
		var claims passport.Claims
		if claims, err = x.IndexX.Verify(ctx, string(ts)); err != nil {
			common.ClearAccessToken(c)
			c.JSON(401, help.E(
				"lark.AuthenticationExpired",
				common.ErrAuthenticationExpired.Error(),
			))
			return
		}

		if _, err = x.LarkX.Link(ctx, claims.UserId, userData); err != nil {
			c.Error(err)
			return
		}
		c.Redirect(302, x.RedirectUrl("/#/authorized", dto.State.Locale))
		return
	}

	var r *LoginResult
	if r, err = x.LarkX.Login(ctx, userData.OpenId); err != nil {
		c.Redirect(302, x.RedirectUrl("/#/unauthorize", dto.State.Locale))
		return
	}

	data := model.NewLogsetLogin(r.User.ID, string(c.GetHeader(x.V.Ip)), "lark", string(c.UserAgent()))
	if err = x.IndexX.WriteLogsetLogin(ctx, data); err != nil {
		c.Error(err)
		return
	}

	common.SetAccessToken(c, r.AccessToken)
	c.Redirect(302, x.RedirectUrl("", dto.State.Locale))
}

func (x *Service) Link(ctx context.Context, userId string, data model.UserLark) (_ *mongo.UpdateResult, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	return x.Db.Collection("users").UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"lark": data}},
	)
}

func (x *Service) Login(ctx context.Context, openId string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.IndexX.Logining(ctx, bson.M{"lark.open_id": openId, "status": true}); err != nil {
		return
	}

	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"lark.open_id": openId, "status": true}).Decode(&r.User); err != nil {
		if err == mongo.ErrNoDocuments {
			err = common.ErrLoginNotExists
			return
		}
		return
	}
	userId := r.User.ID.Hex()
	if r.AccessToken, err = x.IndexX.CreateAccessToken(ctx, userId); err != nil {
		return
	}
	return
}
