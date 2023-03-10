package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/cloudwego/hertz/pkg/common/utils"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/weplanx/openapi/client"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/utils/captcha"
	"github.com/weplanx/utils/locker"
	"github.com/weplanx/utils/passlib"
	"github.com/weplanx/utils/passport"
	"github.com/weplanx/utils/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	*common.Inject
	Passport *passport.Passport
	Locker   *locker.Locker
	Captcha  *captcha.Captcha
	Sessions *sessions.Service
}

func (x *Service) Login(ctx context.Context, email string, password string, logdata *model.LoginLog) (ts string, err error) {
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"email":  email,
			"status": true,
		}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			err = errors.NewPublic("the user does not exist or has been frozen")
			return
		}

		return
	}

	userId := user.ID.Hex()
	logdata.SetUserID(userId)

	var maxLoginFailures bool
	if maxLoginFailures, err = x.Locker.Verify(ctx, userId, x.Values.LoginFailures); err != nil {
		return
	}
	if maxLoginFailures {
		err = errors.NewPublic("the user has exceeded the maximum number of login failures")
		return
	}

	var match bool
	if match, err = passlib.Verify(password, user.Password); err != nil {
		return
	}
	if !match {
		if err = x.Locker.Update(ctx, userId, x.Values.LoginTTL); err != nil {
			return
		}
		err = errors.NewPublic("the user email or password is incorrect")
		return
	}

	jti, _ := gonanoid.Nanoid()
	if ts, err = x.Passport.Create(userId, jti); err != nil {
		return
	}
	if err = x.Locker.Delete(ctx, userId); err != nil {
		return
	}
	if err = x.Sessions.Set(ctx, userId, jti); err != nil {
		return
	}

	key := x.Values.Name("users", userId)
	if _, err = x.Redis.Del(ctx, key).Result(); err != nil {
		return
	}

	return
}

func (x *Service) Openapi() (*client.Client, error) {
	return client.New(
		x.Values.OpenapiUrl,
		client.SetApiGateway(x.Values.OpenapiKey, x.Values.OpenapiSecret),
	)
}

func (x *Service) WriteLoginLog(ctx context.Context, logdata *model.LoginLog) (err error) {
	var oapi *client.Client
	if oapi, err = x.Openapi(); err != nil {
		return
	}
	var detail map[string]interface{}
	if detail, err = oapi.GetIp(ctx, logdata.Metadata.ClientIP); err != nil {
		return
	}
	logdata.SetLocation(detail)

	if _, err = x.Db.Collection("login_logs").InsertOne(ctx, logdata); err != nil {
		return
	}
	return
}

func (x *Service) Verify(ctx context.Context, ts string) (claims passport.Claims, err error) {
	if claims, err = x.Passport.Verify(ts); err != nil {
		return
	}
	var result bool
	if result, err = x.Sessions.Verify(ctx, claims.UserId, claims.ID); err != nil {
		return
	}
	if !result {
		err = errors.NewPublic("the session token is inconsistent")
		return
	}

	// TODO: Check User Status

	if err = x.Sessions.Renew(ctx, claims.UserId); err != nil {
		return
	}

	return
}

func (x *Service) GetRefreshCode(ctx context.Context, userId string) (code string, err error) {
	if code, err = gonanoid.Nanoid(); err != nil {
		return
	}
	if err = x.Captcha.Create(ctx, userId, code, 15*time.Second); err != nil {
		return
	}
	return
}

func (x *Service) RefreshToken(ctx context.Context, claims passport.Claims, code string) (ts string, err error) {
	if err = x.Captcha.Verify(ctx, claims.UserId, code); err != nil {
		return
	}
	if ts, err = x.Passport.Create(claims.UserId, claims.ID); err != nil {
		return
	}
	return
}

func (x *Service) Logout(ctx context.Context, userId string) (err error) {
	return x.Sessions.Remove(ctx, userId)
}

//func (x *Service) GetIdentity(ctx context.Context, userId string) (data model.User, err error) {
//	key := x.Values.Name("users", userId)
//	var exists int64
//	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
//		return
//	}
//
//	if exists == 0 {
//		id, _ := primitive.ObjectIDFromHex(userId)
//		option := options.FindOne().SetProjection(bson.M{"password": 0})
//		if err = x.Db.Collection("users").
//			FindOne(ctx, bson.M{
//				"_id":    id,
//				"status": true,
//			}, option).
//			Decode(&data); err != nil {
//			return
//		}
//
//		var value string
//		if value, err = sonic.MarshalString(data); err != nil {
//			return
//		}
//
//		if err = x.Redis.Set(ctx, key, value, 0).Err(); err != nil {
//			return
//		}
//
//		return
//	}
//
//	var result string
//	if result, err = x.Redis.Get(ctx, key).Result(); err != nil {
//		return
//	}
//	if err = sonic.UnmarshalString(result, &data); err != nil {
//		return
//	}
//
//	return
//}

func (x *Service) GetUser(ctx context.Context, userId string) (data utils.H, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&user); err != nil {
		return
	}

	data = utils.H{
		"email":        user.Email,
		"name":         user.Name,
		"avatar":       user.Avatar,
		"backup_email": user.BackupEmail,
		"status":       user.Status,
		"create_time":  user.CreateTime,
		"update_time":  user.UpdateTime,
	}

	if user.Feishu.OpenId != "" {
		feishu := user.Feishu
		data["feishu"] = utils.H{
			"name":          feishu.Name,
			"en_name":       feishu.EnName,
			"avatar_url":    feishu.AvatarUrl,
			"avatar_thumb":  feishu.AvatarThumb,
			"avatar_middle": feishu.AvatarMiddle,
			"avatar_big":    feishu.AvatarBig,
			"open_id":       feishu.OpenId,
		}
	}

	return
}

func (x *Service) SetUser(ctx context.Context, userId string, update bson.M) (result interface{}, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)

	if result, err = x.Db.Collection("users").
		UpdateByID(ctx, id, update); err != nil {
		return
	}

	key := x.Values.Name("users", userId)
	if _, err = x.Redis.Del(ctx, key).Result(); err != nil {
		return
	}

	return
}
