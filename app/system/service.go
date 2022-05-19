package system

import (
	"api/app/users"
	"api/common"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/weplanx/go/vars"
	openapi "github.com/weplanx/openapi/client"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
	Vars  *vars.Service
	Users *users.Service
}

// GetSessions 获取所有会话
func (x *Service) GetSessions(ctx context.Context) (values []string, err error) {
	var cursor uint64
	for {
		var keys []string
		var next uint64
		if keys, next, err = x.Redis.Scan(ctx,
			cursor, x.Values.KeyName("sessions", "*"), 1000,
		).Result(); err != nil {
			return
		}
		uids := make([]string, len(keys))
		for k, v := range keys {
			uids[k] = strings.Replace(v, x.Values.KeyName("sessions", ""), "", -1)
		}
		values = append(values, uids...)
		if next == 0 {
			break
		}
		cursor = next
	}
	return
}

// VerifySession 验证会话一致性
func (x *Service) VerifySession(ctx context.Context, uid string, jti string) (_ bool, err error) {
	var value string
	if value, err = x.Redis.Get(ctx, x.Values.KeyName("sessions", uid)).Result(); err != nil {
		return
	}
	return value == jti, nil
}

// SetSession 设置会话
func (x *Service) SetSession(ctx context.Context, uid string, jti string) (err error) {
	exp, err := x.Vars.GetUserSessionExpire(ctx)
	if err != nil {
		return
	}
	if err = x.Redis.Set(ctx, x.Values.KeyName("sessions", uid), jti, exp).Err(); err != nil {
		return
	}
	return
}

// RenewSession 续约会话
func (x *Service) RenewSession(ctx context.Context, uid string) (err error) {
	exp, err := x.Vars.GetUserSessionExpire(ctx)
	if err != nil {
		return
	}
	if err = x.Redis.Expire(ctx, x.Values.KeyName("sessions", uid), exp).Err(); err != nil {
		return
	}
	return
}

// DeleteSession 删除会话
func (x *Service) DeleteSession(ctx context.Context, uid string) (err error) {
	return x.Redis.Del(ctx, x.Values.KeyName("sessions", uid)).Err()
}

// DeleteSessions 删除所有会话
func (x *Service) DeleteSessions(ctx context.Context) (err error) {
	var cursor uint64
	var keys []string
	for {
		var next uint64
		if keys, next, err = x.Redis.Scan(ctx,
			cursor, x.Values.KeyName("sessions", "*"), 1000,
		).Result(); err != nil {
			return
		}
		if next == 0 {
			break
		}
		cursor = next
	}
	return x.Redis.Del(ctx, keys...).Err()
}

// CheckLockForUser 判断用户是否锁定
func (x *Service) CheckLockForUser(ctx context.Context, uid string) (err error) {
	key := x.Values.KeyName("lock", uid)
	var count int64
	if count, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}
	if count == 0 {
		return
	}
	times, err := x.Redis.Get(ctx, key).Int()
	if err != nil {
		return
	}
	userLoginFailedTimes, err := x.Vars.GetUserLoginFailedTimes(ctx)
	if err != nil {
		return
	}
	userLockTime, err := x.Vars.GetUserLockTime(ctx)
	if err != nil {
		return
	}
	// 用户连续登录失败已超出最大次数
	if times > userLoginFailedTimes {
		// 针对锁定缓存延长锁定时效
		if err = x.Redis.Expire(ctx, key, userLockTime).Err(); err != nil {
			return
		}
		return errors.New("用户连续登录失败已超出最大次数")
	}
	return
}

// IncLockForUser 增加锁定次数
func (x *Service) IncLockForUser(ctx context.Context, uid string) (err error) {
	key := x.Values.KeyName("lock", uid)
	if err = x.Redis.Incr(ctx, key).Err(); err != nil {
		return
	}
	return
}

// CreateVerifyCode 创建验证码
func (x *Service) CreateVerifyCode(ctx context.Context, name string, code string, ttl time.Duration) error {
	return x.Redis.Set(ctx, x.Values.KeyName("verify", name), code, ttl).Err()
}

// ExistsVerifyCode 已存在的验证码
func (x *Service) ExistsVerifyCode(ctx context.Context, name string) (exists bool, err error) {
	var count int64
	if count, err = x.Redis.Exists(ctx, x.Values.KeyName("verify", name)).Result(); err != nil {
		return
	}
	return count != 0, nil
}

// CheckVerifyCode 校验验证码
func (x *Service) CheckVerifyCode(ctx context.Context, name string, code string) (result bool, err error) {
	var value string
	if value, err = x.Redis.Get(ctx, x.Values.KeyName("verify", name)).Result(); err != nil {
		return
	}
	return value == code, nil
}

// DeleteVerifyCode 移除验证码
func (x *Service) DeleteVerifyCode(ctx context.Context, name string) error {
	return x.Redis.Del(ctx, x.Values.KeyName("verify", name)).Err()
}

// OpenAPI 开放服务客户端
func (x *Service) OpenAPI(ctx context.Context) (_ *openapi.OpenAPI, err error) {
	url, err := x.Vars.GetOpenapiUrl(ctx)
	if err != nil {
		return
	}
	key, err := x.Vars.GetOpenapiKey(ctx)
	if err != nil {
		return
	}
	secret, err := x.Vars.GetOpenapiSecret(ctx)
	if err != nil {
		return
	}
	return openapi.New(url, openapi.SetCertification(key, secret)), nil
}

// PushLoginLog 推送登录日志
func (x *Service) PushLoginLog(ctx context.Context, doc *common.LoginLogDto) (err error) {
	var client *openapi.OpenAPI
	if client, err = x.OpenAPI(ctx); err != nil {
		return
	}
	if doc.Detail, err = client.Ip(ctx, doc.Ip); err != nil {
		return
	}
	if err = x.Users.UpdateOneById(ctx, doc.User, bson.M{
		"$inc": bson.M{"sessions": 1},
		"$set": bson.M{
			"last": fmt.Sprintf(`%s %s`, doc.Detail["isp"], doc.Ip),
		},
	}); err != nil {
		return err
	}
	if _, err = x.Db.Collection("login_logs").InsertOne(ctx, doc); err != nil {
		return
	}
	return
}

func (x *Service) SendEmail(ctx context.Context, to []string, name string, subject string, html []byte) (err error) {
	host, err := x.Vars.GetEmailHost(ctx)
	if err != nil {
		return
	}
	port, err := x.Vars.GetEmailPort(ctx)
	if err != nil {
		return
	}
	username, err := x.Vars.GetEmailUsername(ctx)
	if err != nil {
		return
	}
	password, err := x.Vars.GetEmailPassword(ctx)
	if err != nil {
		return
	}
	e := &email.Email{
		To:      to,
		From:    fmt.Sprintf(`%s <%s>`, name, username),
		Subject: subject,
		HTML:    html,
	}
	if err = e.SendWithTLS(
		fmt.Sprintf(`%s:%s`, host, port),
		smtp.PlainAuth("",
			username,
			password,
			host,
		),
		&tls.Config{
			ServerName: host,
		},
	); err != nil {
		panic(err)
	}
	return
}

type EmailVerifyDto struct {
	Name string
	User string
	Code string
	Year int
}

// EmailCode 邮箱验证码
func (x *Service) EmailCode(ctx context.Context, user string, code string, to []string) (err error) {
	var tpl *template.Template
	if tpl, err = template.ParseFiles("./templates/email_verify.gohtml"); err != nil {
		return
	}
	dto := EmailVerifyDto{
		Name: x.Values.Name,
		User: user,
		Code: code,
		Year: time.Now().Year(),
	}
	var buf bytes.Buffer
	if err = tpl.Execute(&buf, dto); err != nil {
		return
	}
	return x.SendEmail(ctx, to, dto.Name, "用户密码重置验证", buf.Bytes())
}
