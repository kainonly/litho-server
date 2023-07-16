package index

import (
	"context"
	"github.com/google/uuid"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passlib"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	*common.Inject
	Sessions *sessions.Service
}

func (x *Service) Login(ctx context.Context, email string, password string) (ts string, err error) {
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"email": email, "status": true}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			err = common.ErrLoginNotExists
			return
		}
		return
	}

	userId := user.ID.Hex()
	if err = x.Locker.Verify(ctx, userId, x.V.LoginFailures); err != nil {
		if err == locker.ErrLocked {
			err = common.ErrLoginMaxFailures
			return
		}
		return
	}

	if err = passlib.Verify(password, user.Password); err != nil {
		if err == passlib.ErrNotMatch {
			x.Locker.Update(ctx, userId, x.V.LoginTTL)
			err = common.ErrLoginInvalid
			return
		}
		return
	}

	jti := uuid.New().String()
	if ts, err = x.Passport.Create(userId, jti); err != nil {
		return
	}
	if status := x.Sessions.Set(ctx, userId, jti); status != "OK" {
		err = common.ErrSession
		return
	}
	x.Locker.Delete(ctx, userId)

	// Refresh user cache
	key := x.V.Name("users", userId)
	if _, err = x.RDb.Del(ctx, key).Result(); err != nil {
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
