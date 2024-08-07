package index

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/jordan-wright/email"
	"github.com/weplanx/go/help"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/smtp"
	"time"
)

type GetForgetCodeDto struct {
	Email string `query:"email" vd:"email"`
}

func (x *Controller) GetForgetCode(ctx context.Context, c *app.RequestContext) {
	var dto GetForgetCodeDto
	if err := c.BindAndValidate(&dto); err != nil {
		c.Error(err)
		return
	}

	if err := x.IndexX.GetForgetCode(ctx, dto.Email); err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
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
