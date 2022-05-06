package system

import (
	"api/app/users"
	"api/common"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
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
	expiration := x.GetExpiration(ctx)
	if err = x.Redis.Set(ctx, x.Values.KeyName("sessions", uid), jti, expiration).Err(); err != nil {
		return
	}
	return
}

// RenewSession 续约会话
func (x *Service) RenewSession(ctx context.Context, uid string) (err error) {
	expiration := x.GetExpiration(ctx)
	if err = x.Redis.Expire(ctx,
		x.Values.KeyName("sessions", uid), expiration,
	).Err(); err != nil {
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

// GetVar 获取变量
func (x *Service) GetVar(ctx context.Context, key string) (value string, err error) {
	if err = x.RefreshVars(ctx); err != nil {
		return
	}
	return x.Redis.HGet(ctx, x.Values.KeyName("vars"), key).Result()
}

// GetVars 获取指定变量
func (x *Service) GetVars(ctx context.Context, keys []string) (values map[string]interface{}, err error) {
	if err = x.RefreshVars(ctx); err != nil {
		return
	}
	var result []interface{}
	if result, err = x.Redis.HMGet(ctx, x.Values.KeyName("vars"), keys...).Result(); err != nil {
		return
	}
	values = make(map[string]interface{})
	for k, v := range keys {
		values[v] = result[k]
	}
	return
}

// RefreshVars 刷新变量
func (x *Service) RefreshVars(ctx context.Context) (err error) {
	key := x.Values.KeyName("vars")
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var cursor *mongo.Cursor
		if cursor, err = x.Db.Collection("vars").Find(ctx, bson.M{}); err != nil {
			return
		}
		var data []common.Var
		if err = cursor.All(ctx, &data); err != nil {
			return
		}
		pipe := x.Redis.Pipeline()
		for _, v := range data {
			switch x := v.Value.(type) {
			case primitive.A:
				b, _ := jsoniter.Marshal(x)
				pipe.HSet(ctx, key, v.Key, b)
				break
			case primitive.M:
				b, _ := jsoniter.Marshal(x)
				pipe.HSet(ctx, key, v.Key, b)
				break
			default:
				pipe.HSet(ctx, key, v.Key, x)
			}
		}
		if _, err = pipe.Exec(ctx); err != nil {
			return
		}
	}
	return
}

// SetVar 设置变量
func (x *Service) SetVar(ctx context.Context, key string, value interface{}) (err error) {
	var exists int64
	if exists, err = x.Db.Collection("vars").CountDocuments(ctx, bson.M{"key": key}); err != nil {
		return
	}
	doc := common.NewVar(key, value)
	if exists == 0 {
		if _, err = x.Db.Collection("vars").InsertOne(ctx, doc); err != nil {
			return
		}
	} else {
		if _, err = x.Db.Collection("vars").ReplaceOne(ctx, bson.M{"key": key}, doc); err != nil {
			return
		}
	}
	if err = x.Redis.Del(ctx, x.Values.KeyName("vars")).Err(); err != nil {
		return
	}
	return
}

// GetExpiration 获取会话有效时间
func (x *Service) GetExpiration(ctx context.Context) (t time.Duration) {
	value, _ := x.GetVar(ctx, "user_session_expire")
	if value != "" {
		t, _ = time.ParseDuration(value)
	} else {
		t = time.Hour
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

// PushLoginLog 推送登录日志
func (x *Service) PushLoginLog(ctx context.Context, doc *common.LoginLogDto) (err error) {
	if doc.Detail, err = x.Open.Ip(ctx, doc.Ip); err != nil {
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
	var option map[string]interface{}
	if option, err = x.GetVars(ctx, []string{
		"email_host",
		"email_port",
		"email_username",
		"email_password",
	}); err != nil {
		return
	}
	e := &email.Email{
		To:      to,
		From:    fmt.Sprintf(`%s <%s>`, name, option["email_username"]),
		Subject: subject,
		HTML:    html,
	}
	if err = e.SendWithTLS(
		fmt.Sprintf(`%s:%s`, option["email_host"], option["email_port"]),
		smtp.PlainAuth("",
			option["email_username"].(string),
			option["email_password"].(string),
			option["email_host"].(string),
		),
		&tls.Config{
			ServerName: option["email_host"].(string),
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
