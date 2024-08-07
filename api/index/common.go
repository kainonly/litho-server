package index

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/google/wire"
	"github.com/weplanx/go/csrf"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/server/api/tencent"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	V    *common.Values
	Csrf *csrf.Csrf

	IndexX *Service
}

type Service struct {
	*common.Inject

	Sessions *sessions.Service
	Passport *common.APIPassport
	TencentX *tencent.Service
}

type M = map[string]interface{}

func R(code string, msg string) M {
	return M{
		"code": code,
		"msg":  msg,
	}
}

type LoginResult struct {
	User        model.User
	AccessToken string
}

func (x *Service) Logining(ctx context.Context, filter bson.M) (u model.User, err error) {
	if err = x.Db.Collection("users").FindOne(ctx, filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			err = common.ErrLoginNotExists
			return
		}
		return
	}

	if err = x.Locker.Verify(ctx, u.ID.Hex(), x.V.LoginFailures); err != nil {
		switch err {
		case locker.ErrLockerNotExists:
			err = nil
			break
		case locker.ErrLocked:
			err = common.ErrLoginMaxFailures
			return
		default:
			return
		}
	}

	return
}

func (x *Service) CreateAccessToken(ctx context.Context, userId string) (ts string, err error) {
	jti := help.Uuid()
	if ts, err = x.Passport.Create(userId, jti, time.Hour*2); err != nil {
		return
	}
	if status := x.Sessions.Set(ctx, userId, jti); status != "OK" {
		err = common.ErrSession
		return
	}

	x.Locker.Delete(ctx, userId)
	return
}

func (x *Service) WriteLogsetLogin(ctx context.Context, data *model.LogsetLogin) (err error) {
	ip := net.ParseIP(data.Metadata.ClientIP)
	if ip == nil {
		return
	}
	var r tencent.IpResult
	if ip.To4() != nil {
		if r, err = x.TencentX.GetIpv4(ctx, data.Metadata.ClientIP); err != nil {
			return
		}
	} else {
		if r, err = x.TencentX.GetIpv6(ctx, data.Metadata.ClientIP); err != nil {
			return
		}
	}
	if !r.(tencent.IpResult).IsSuccess() {
		return errors.NewPublic(r.GetMsg())
	}

	data.SetVersion("shuliancloud.v4")
	data.SetDetail(r.GetDetail())
	if _, err = x.Db.Collection("logset_logins").InsertOne(ctx, data); err != nil {
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
