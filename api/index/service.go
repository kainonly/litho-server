package index

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base32"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/jordan-wright/email"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passlib"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/go/totp"
	"github.com/weplanx/server/api/tencent"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net"
	"net/smtp"
	"net/url"
	"time"
)

type Service struct {
	*common.Inject
	Sessions *sessions.Service
	Tencent  *tencent.Service
	Passport *common.APIPassport
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

func (x *Service) Login(ctx context.Context, username string, password string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.Logining(ctx, bson.M{"email": username, "status": true}); err != nil {
		return
	}

	userId := r.User.ID.Hex()
	if err = passlib.Verify(password, r.User.Password); err != nil {
		if err == passlib.ErrNotMatch {
			x.Locker.Update(ctx, userId, x.V.LoginTTL)
			err = common.ErrLoginInvalid
			return
		}
		return
	}

	if r.AccessToken, err = x.CreateAccessToken(ctx, userId); err != nil {
		return
	}
	return
}

func (x *Service) GetLoginSms(ctx context.Context, phone string) (code string, err error) {
	keyLock := fmt.Sprintf(`phone:%s`, phone)
	if err = x.Locker.Verify(ctx, keyLock, x.V.LoginFailures); err != nil {
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

	var n int64
	if n, err = x.Db.Collection("users").
		CountDocuments(ctx, bson.M{"phone": phone, "status": true}); err != nil {
		return
	}
	if n == 0 {
		x.Locker.Update(ctx, keyLock, time.Hour*24)
		err = common.ErrSmsNotExists
		return
	}

	key := fmt.Sprintf(`sms-login:%s`, phone)
	if exists := x.Captcha.Exists(ctx, key); exists {
		err = common.ErrCodeFrequently
		return
	}

	code = help.RandomNumber(6)
	// TODO: Change to values...
	if err = x.Tencent.SmsSend(ctx, "WEB应用技术分享网", "1889615", []string{code}, []string{phone}); err != nil {
		return
	}

	x.Captcha.Create(ctx, key, code, time.Minute*2)
	return
}

func (x *Service) LoginSms(ctx context.Context, phone string, code string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.Logining(ctx, bson.M{"phone": phone, "status": true}); err != nil {
		return
	}

	userId := r.User.ID.Hex()
	key := fmt.Sprintf(`sms-login:%s`, phone)
	if err = x.Captcha.Verify(ctx, key, code); err != nil {
		x.Locker.Update(ctx, userId, x.V.LoginTTL)
		err = common.ErrSmsInvalid
		return
	}

	if r.AccessToken, err = x.CreateAccessToken(ctx, userId); err != nil {
		return
	}
	x.Captcha.Delete(ctx, key)
	return
}

func (x *Service) LoginTotp(ctx context.Context, email string, code string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.Logining(ctx, bson.M{"email": email, "status": true}); err != nil {
		return
	}

	userId := r.User.ID.Hex()
	otpc := &totp.Totp{
		Secret:  r.User.Totp,
		Window:  1,
		Counter: 0,
	}
	var check bool
	if check, err = otpc.Authenticate(code); err != nil {
		return
	}
	if !check {
		x.Locker.Update(ctx, userId, x.V.LoginTTL)
		err = common.ErrLoginInvalid
		return
	}

	if r.AccessToken, err = x.CreateAccessToken(ctx, userId); err != nil {
		return
	}
	return
}

func (x *Service) WriteLogsetLogin(ctx context.Context, data *model.LogsetLogin) (err error) {
	ip := net.ParseIP(data.Metadata.ClientIP)
	if ip == nil {
		return
	}
	var r tencent.IpResult
	if ip.To4() != nil {
		if r, err = x.Tencent.GetIpv4(ctx, data.Metadata.ClientIP); err != nil {
			return
		}
	} else {
		if r, err = x.Tencent.GetIpv6(ctx, data.Metadata.ClientIP); err != nil {
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

func (x *Service) GetForgetCode(ctx context.Context, username string) (err error) {
	keyLock := fmt.Sprintf(`forget:%s`, username)
	if err = x.Locker.Verify(ctx, keyLock, x.V.LoginFailures); err != nil {
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

	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"email": username, "status": true}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			x.Locker.Update(ctx, keyLock, time.Hour*24)
			err = common.ErrEmailNotExists
			return
		}
	}

	key := fmt.Sprintf(`forget:%s`, username)
	if exists := x.Captcha.Exists(ctx, key); exists {
		err = common.ErrCodeFrequently
		return
	}

	code := help.RandomNumber(6)
	var tpl *template.Template
	if tpl, err = template.ParseFiles("./templates/email_code.gohtml"); err != nil {
		return
	}

	var buf bytes.Buffer
	name := fmt.Sprintf("用户 <%s>", username)
	if user.Name != "" {
		name = user.Name
	}
	if err = tpl.Execute(&buf, M{
		"Name": name,
		"Code": code,
		"Year": time.Now().Year(),
	}); err != nil {
		return
	}

	mail := &email.Email{
		To:      []string{username},
		From:    fmt.Sprintf(`WEPLANX <%s>`, x.V.EmailUsername),
		Subject: "邮箱验证",
		HTML:    buf.Bytes(),
	}
	if err = mail.SendWithTLS(
		fmt.Sprintf(`%s:%d`, x.V.EmailHost, x.V.EmailPort),
		smtp.PlainAuth(
			"",
			x.V.EmailUsername,
			x.V.EmailPassword,
			x.V.EmailHost,
		),
		&tls.Config{ServerName: x.V.EmailHost},
	); err != nil {
		return
	}

	x.Captcha.Create(ctx, key, code, time.Minute*15)
	return
}

func (x *Service) ForgetReset(ctx context.Context, dto ForgetResetDto) (err error) {
	key := fmt.Sprintf(`forget:%s`, dto.Email)
	if err = x.Captcha.Verify(ctx, key, dto.Code); err != nil {
		err = common.ErrEmailInvalid
		return
	}

	filter := bson.M{"email": dto.Email, "status": true}
	password, _ := passlib.Hash(dto.Password)
	data := bson.M{"$set": bson.M{"password": password}}
	if _, err = x.Db.Collection("users").UpdateOne(ctx, filter, data); err != nil {
		return
	}

	x.Captcha.Delete(ctx, key)
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

	// TODO: Check user status

	x.Sessions.Renew(ctx, claims.UserId)

	return
}

func (x *Service) GetRefreshCode(ctx context.Context, userId string) (code string, err error) {
	code = help.Uuid()
	x.Captcha.Create(ctx, userId, code, 15*time.Second)
	return
}

func (x *Service) RefreshToken(ctx context.Context, claims passport.Claims, code string) (ts string, err error) {
	if err = x.Captcha.Verify(ctx, claims.UserId, code); err != nil {
		return
	}
	if ts, err = x.Passport.Create(claims.UserId, claims.ID, time.Hour*2); err != nil {
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

func (x *Service) GetUserPhoneCode(ctx context.Context, phone string) (code string, err error) {
	code = help.RandomNumber(6)
	// TODO: Change to values
	if err = x.Tencent.SmsSend(ctx, "WEB应用技术分享网", "1889620", []string{code}, []string{phone}); err != nil {
		return
	}
	key := fmt.Sprintf(`sms-bind:%s`, phone)
	if exists := x.Captcha.Exists(ctx, key); exists {
		err = common.ErrCodeFrequently
		return
	}
	x.Captcha.Create(ctx, key, code, time.Minute*2)
	return
}

func (x *Service) SetUserPhone(ctx context.Context, userId string, phone string, code string) (r interface{}, err error) {
	key := fmt.Sprintf(`sms-bind:%s`, phone)
	if err = x.Captcha.Verify(ctx, key, code); err != nil {
		err = common.ErrSmsInvalid
		return
	}

	x.Captcha.Delete(ctx, key)
	return x.SetUser(ctx, userId, bson.M{
		"$set": bson.M{
			"phone": phone,
		},
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
