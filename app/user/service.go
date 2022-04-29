package user

import (
	"api/app/users"
	"api/app/vars"
	"api/common"
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/smtp"
	"time"
)

type Service struct {
	*common.Inject
	Vars  *vars.Service
	Users *users.Service
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
