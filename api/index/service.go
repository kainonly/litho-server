package index

import (
	"github.com/weplanx/server/common"
)

type Service struct {
	*common.Inject
}

//func (x *Service) Login(ctx context.Context, email string, password string) (ts string, err error) {
//	var user model.User
//	if err = x.Db.Collection("users").
//		FindOne(ctx, bson.M{
//			"email":  email,
//			"status": true,
//		}).Decode(&user); err != nil {
//		if err == mongo.ErrNoDocuments {
//			err = errors.NewPublic("the user does not exist or has been frozen")
//			return
//		}
//
//		return
//	}
//
//	userId := user.ID.Hex()
//
//	var maxLoginFailures bool
//	if maxLoginFailures, err = x.Locker.Verify(ctx, userId, x.V.LoginFailures); err != nil {
//		return
//	}
//	if maxLoginFailures {
//		err = errors.NewPublic("the user has exceeded the maximum number of login failures")
//		return
//	}
//
//	var match bool
//	if match, err = passlib.Verify(password, user.Password); err != nil {
//		return
//	}
//	if !match {
//		if err = x.Locker.Update(ctx, userId, x.V.LoginTTL); err != nil {
//			return
//		}
//		err = errors.NewPublic("the user email or password is incorrect")
//		return
//	}
//
//	jti, _ := gonanoid.Nanoid()
//	if ts, err = x.Passport.Create(userId, jti); err != nil {
//		return
//	}
//	if err = x.Locker.Delete(ctx, userId); err != nil {
//		return
//	}
//	if err = x.Sessions.Set(ctx, userId, jti); err != nil {
//		return
//	}
//
//	key := x.V.Name("users", userId)
//	if _, err = x.RDb.Del(ctx, key).Result(); err != nil {
//		return
//	}
//
//	return
//}
