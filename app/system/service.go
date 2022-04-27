package system

import (
	"api/app/users"
	"api/common"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jordan-wright/email"
	jsoniter "github.com/json-iterator/go"
	"github.com/weplanx/go/helper"
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

func (x *Service) AppName() string {
	return x.Values.Namespace
}

// GetVars 获取指定变量
func (x *Service) GetVars(ctx context.Context, keys []string) (data map[string]interface{}, err error) {
	if err = x.RefreshVars(ctx); err != nil {
		return
	}
	var values []interface{}
	if values, err = x.Redis.HMGet(ctx, x.Values.KeyName("vars"), keys...).Result(); err != nil {
		return
	}
	data = make(map[string]interface{})
	for k, v := range keys {
		data[v] = values[k]
	}
	return
}

// GetVar 获取变量
func (x *Service) GetVar(ctx context.Context, key string) (value string, err error) {
	if err = x.RefreshVars(ctx); err != nil {
		return
	}
	return x.Redis.HGet(ctx, x.Values.KeyName("vars"), key).Result()
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
		var data []map[string]interface{}
		if err = cursor.All(ctx, &data); err != nil {
			return
		}
		pipe := x.Redis.Pipeline()
		for _, v := range data {
			switch x := v["value"].(type) {
			case primitive.A:
				b, _ := jsoniter.Marshal(x)
				pipe.HSet(ctx, key, v["key"], b)
				break
			case primitive.M:
				b, _ := jsoniter.Marshal(x)
				pipe.HSet(ctx, key, v["key"], b)
				break
			default:
				pipe.HSet(ctx, key, v["key"], x)
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
	doc := bson.M{"key": key, "value": value}
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

// GetSessions 获取所有会话
func (x *Service) GetSessions(ctx context.Context) (values []string, err error) {
	var cursor uint64
	for {
		var keys []string
		var next uint64
		if keys, next, err = x.Redis.Scan(ctx,
			cursor, x.Values.KeyName("session", "*"), 1000,
		).Result(); err != nil {
			return
		}
		uids := make([]string, len(keys))
		for k, v := range keys {
			uids[k] = strings.Replace(v, x.Values.KeyName("session", ""), "", -1)
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
	if value, err = x.Redis.Get(ctx, x.Values.KeyName("session", uid)).Result(); err != nil {
		return
	}
	return value == jti, nil
}

// GetExpiration 获取会话有效时间
func (x *Service) GetExpiration(ctx context.Context) (t time.Duration, err error) {
	value, _ := x.GetVar(ctx, "user_session_expire")
	t = time.Hour
	if value != "" {
		if t, err = time.ParseDuration(value); err != nil {
			return
		}
	}
	return
}

// SetSession 设置会话
func (x *Service) SetSession(ctx context.Context, uid string, jti string) (err error) {
	var expiration time.Duration
	if expiration, err = x.GetExpiration(ctx); err != nil {
		return
	}
	if err = x.Redis.Set(ctx,
		x.Values.KeyName("session", uid), jti, expiration,
	).Err(); err != nil {
		return
	}
	return
}

// RenewSession 续约会话
func (x *Service) RenewSession(ctx context.Context, uid string) (err error) {
	var expiration time.Duration
	if expiration, err = x.GetExpiration(ctx); err != nil {
		return
	}
	if err = x.Redis.Expire(ctx,
		x.Values.KeyName("session", uid), expiration,
	).Err(); err != nil {
		return
	}
	return
}

// DeleteSessions 删除所有会话
func (x *Service) DeleteSessions(ctx context.Context) (err error) {
	var cursor uint64
	var keys []string
	for {
		var next uint64
		if keys, next, err = x.Redis.Scan(ctx,
			cursor, x.Values.KeyName("session", "*"), 1000,
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

// DeleteSession 删除会话
func (x *Service) DeleteSession(ctx context.Context, uid string) (err error) {
	return x.Redis.Del(ctx, x.Values.KeyName("session", uid)).Err()
}

// CreateCode 创建验证码
func (x *Service) CreateCode(ctx context.Context, name string, code string, ttl time.Duration) error {
	return x.Redis.Set(ctx, x.Values.KeyName("verify", name), code, ttl).Err()
}

func (x *Service) ExistsCode(ctx context.Context, name string) (exists bool, err error) {
	var count int64
	if count, err = x.Redis.Exists(ctx, x.Values.KeyName("verify", name)).Result(); err != nil {
		return
	}
	return count != 0, nil
}

// VerifyCode 校验验证码
func (x *Service) VerifyCode(ctx context.Context, name string, code string) (result bool, err error) {
	var value string
	if value, err = x.Redis.Get(ctx, x.Values.KeyName("verify", name)).Result(); err != nil {
		return
	}
	return value == code, nil
}

// DeleteCode 移除验证码
func (x *Service) DeleteCode(ctx context.Context, name string) error {
	return x.Redis.Del(ctx, x.Values.KeyName("verify", name)).Err()
}

type EmailVerifyDto struct {
	Name string
	User string
	Code string
	Year int
}

// EmailCode 邮箱验证码
func (x *Service) EmailCode(user string, code string, to []string) (err error) {
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
	option := x.Values.Email
	e := &email.Email{
		To:      to,
		From:    fmt.Sprintf(`%s <%s>`, dto.Name, option.Username),
		Subject: "用户密码重置验证",
		HTML:    buf.Bytes(),
	}
	if err = e.SendWithTLS(
		fmt.Sprintf(`%s:%s`, option.Host, option.Port),
		smtp.PlainAuth("", option.Username, option.Password, option.Host),
		&tls.Config{
			ServerName: option.Host,
		},
	); err != nil {
		panic(err)
	}
	return
}

func (x *Service) WriteLoginLog(ctx context.Context, doc *common.LoginLogDto) (err error) {
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

func (x *Service) Sort(ctx context.Context, model string, sort []primitive.ObjectID) (*mongo.BulkWriteResult, error) {
	var models []mongo.WriteModel
	for i, oid := range sort {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": oid}).
			SetUpdate(bson.M{"$set": bson.M{"sort": i}}),
		)
	}
	return x.Db.Collection(model).BulkWrite(ctx, models)
}

// Uploader 上传预签名
func (x *Service) Uploader() (data interface{}, err error) {
	option := x.Values.QCloud
	expired := time.Second * time.Duration(option.Cos.Expired)
	date := time.Now()
	keyTime := fmt.Sprintf(`%d;%d`, date.Unix(), date.Add(expired).Unix())
	key := fmt.Sprintf(`%s/%s/%s`,
		x.AppName(),
		date.Format("20060102"),
		helper.Uuid(),
	)
	policy := map[string]interface{}{
		"expiration": date.Add(expired).Format("2006-01-02T15:04:05.000Z"),
		"conditions": []interface{}{
			map[string]interface{}{"bucket": option.Cos.Bucket},
			[]interface{}{"starts-with", "$key", key},
			map[string]interface{}{"q-sign-algorithm": "sha1"},
			map[string]interface{}{"q-ak": option.SecretID},
			map[string]interface{}{"q-sign-time": keyTime},
		},
	}
	var policyText []byte
	if policyText, err = jsoniter.Marshal(policy); err != nil {
		return
	}
	signKeyHash := hmac.New(sha1.New, []byte(option.SecretKey))
	signKeyHash.Write([]byte(keyTime))
	signKey := hex.EncodeToString(signKeyHash.Sum(nil))
	stringToSignHash := sha1.New()
	stringToSignHash.Write(policyText)
	stringToSign := hex.EncodeToString(stringToSignHash.Sum(nil))
	signatureHash := hmac.New(sha1.New, []byte(signKey))
	signatureHash.Write([]byte(stringToSign))
	signature := hex.EncodeToString(signatureHash.Sum(nil))
	return gin.H{
		"key":              key,
		"policy":           policyText,
		"q-sign-algorithm": "sha1",
		"q-ak":             option.SecretID,
		"q-key-time":       keyTime,
		"q-signature":      signature,
	}, nil
}
