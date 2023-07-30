package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/google/uuid"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passlib"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/server/api/tencent"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	*common.Inject
	Sessions *sessions.Service
	Tencent  *tencent.Service
}

type LoginResult struct {
	User        model.User
	AccessToken string
}

func (x *Service) Login(ctx context.Context, email string, password string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"email": email, "status": true}).Decode(&r.User); err != nil {
		if err == mongo.ErrNoDocuments {
			err = common.ErrLoginNotExists
			return
		}
		return
	}

	userId := r.User.ID.Hex()
	if err = x.Locker.Verify(ctx, userId, x.V.LoginFailures); err != nil {
		switch err {
		case locker.ErrLockerNotExists:
			break
		case locker.ErrLocked:
			err = common.ErrLoginMaxFailures
			return
		default:
			return
		}
	}

	if err = passlib.Verify(password, r.User.Password); err != nil {
		if err == passlib.ErrNotMatch {
			x.Locker.Update(ctx, userId, x.V.LoginTTL)
			err = common.ErrLoginInvalid
			return
		}
		return
	}

	jti := uuid.New().String()
	if r.AccessToken, err = x.Passport.Create(userId, jti); err != nil {
		return
	}
	if status := x.Sessions.Set(ctx, userId, jti); status != "OK" {
		err = common.ErrSession
		return
	}
	x.Locker.Delete(ctx, userId)

	// Refresh user cache
	key := x.V.Name("users", userId)
	if err = x.RDb.Del(ctx, key).Err(); err != nil {
		return
	}

	return
}

func (x *Service) WriteLogsetLogined(ctx context.Context, data *model.LogsetLogined) (err error) {
	var r *tencent.CityResult
	if r, err = x.Tencent.GetCity(ctx, data.Metadata.ClientIP); err != nil {
		return
	}
	if !r.Success {
		return errors.NewPublic(r.Msg)
	}
	data.SetVersion("shuliancloud.v4")
	data.SetDetail(r.GetDetail())
	if _, err = x.Db.Collection("logset_logined").InsertOne(ctx, data); err != nil {
		return
	}
	filter := bson.M{"_id": data.Metadata.UserID}
	if _, err = x.Db.Collection("users").UpdateOne(ctx, filter, bson.M{
		"$inc": bson.M{"sessions": 1},
		"$set": bson.M{"history": data},
	}); err != nil {
		return
	}
	return
}

func (x *Service) Verify(ctx context.Context, ts string) (claims passport.Claims, err error) {
	if claims, err = x.Passport.Verify(ts); err != nil {
		return
	}
	result := x.Sessions.Verify(ctx, claims.UserId, claims.ID)
	if !result {
		err = common.ErrSessionInconsistent
		return
	}

	// TODO: Check User Status

	x.Sessions.Renew(ctx, claims.UserId)

	return
}

func (x *Service) GetRefreshCode(ctx context.Context, userId string) (code string, err error) {
	code = uuid.New().String()
	x.Captcha.Create(ctx, userId, code, 15*time.Second)
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

func (x *Service) Logout(ctx context.Context, userId string) {
	x.Sessions.Remove(ctx, userId)
}

func (x *Service) GetUser(ctx context.Context, userId string) (data M, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"_id": id}).
		Decode(&user); err != nil {
		return
	}

	data = M{
		"email":        user.Email,
		"name":         user.Name,
		"avatar":       user.Avatar,
		"backup_email": user.BackupEmail,
		"sessions":     user.Sessions,
		"history":      user.History,
		"status":       user.Status,
		"create_time":  user.CreateTime,
		"update_time":  user.UpdateTime,
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

func (x *Service) SetUser(ctx context.Context, userId string, update bson.M) (result interface{}, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)

	if result, err = x.Db.Collection("users").
		UpdateByID(ctx, id, update); err != nil {
		return
	}

	key := x.V.Name("users", userId)
	if _, err = x.RDb.Del(ctx, key).Result(); err != nil {
		return
	}

	return
}
