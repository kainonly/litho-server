package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (x *Controller) GetUser(ctx context.Context, c *app.RequestContext) {
	claims := common.Claims(c)
	data, err := x.IndexX.GetUser(ctx, claims.UserId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, data)
}

func (x *Service) GetUser(ctx context.Context, userId string) (data M, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&user); err != nil {
		return
	}

	detail := M{}
	if user.History != nil {
		for _, v := range user.History.Detail.(bson.D) {
			detail[v.Key] = v.Value
		}
		user.History.Detail = detail
	}

	phoneStatus := ""
	if user.Phone != "" {
		phoneStatus = "*"
	}

	totpStatus := ""
	if user.Totp != "" {
		totpStatus = "*"
	}

	data = M{
		"_id":         user.ID,
		"email":       user.Email,
		"name":        user.Name,
		"avatar":      user.Avatar,
		"phone":       phoneStatus,
		"sessions":    user.Sessions,
		"history":     user.History,
		"totp":        totpStatus,
		"status":      user.Status,
		"create_time": user.CreateTime,
		"update_time": user.UpdateTime,
	}

	if user.Lark != nil {
		lark := user.Lark
		data["lark"] = M{
			"name":          lark.Name,
			"en_name":       lark.EnName,
			"avatar_url":    lark.AvatarUrl,
			"avatar_thumb":  lark.AvatarThumb,
			"avatar_middle": lark.AvatarMiddle,
			"avatar_big":    lark.AvatarBig,
			"open_id":       lark.OpenId,
		}
	}

	return
}
