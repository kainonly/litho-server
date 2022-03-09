package index

import (
	"api/common"
	"api/model"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/weplanx/go/helper"
	"github.com/weplanx/go/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) AppName() string {
	return x.Values.Name
}

func (x *Service) Install(ctx context.Context, value InstallDto) (err error) {
	// 初始化权限组
	role := model.NewRole("超级管理员").
		SetDescription("系统默认设置").
		SetLabel("最高权限")
	var result *mongo.InsertOneResult
	if result, err = x.Db.Collection("roles").
		InsertOne(ctx, role); err != nil {
		return
	}
	if _, err = x.Db.Collection("roles").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"name": 1},
				Options: options.Index().SetName("uk_name").SetUnique(true),
			},
			{
				Keys:    bson.M{"labels": 1},
				Options: options.Index().SetName("idx_labels"),
			},
		},
	); err != nil {
		return
	}
	// 初始化管理用户
	var pwd string
	if pwd, err = password.Create(value.Password); err != nil {
		return
	}
	user := model.NewUser("kain", pwd).
		AddEmail(value.Email).
		SetRoles([]primitive.ObjectID{result.InsertedID.(primitive.ObjectID)})
	if _, err = x.Db.Collection("users").
		InsertOne(ctx, user); err != nil {
		return
	}
	if _, err = x.Db.Collection("users").Indexes().CreateMany(ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"username": 1},
				Options: options.Index().SetName("uk_username").SetUnique(true),
			},
			{
				Keys:    bson.M{"labels": 1},
				Options: options.Index().SetName("idx_labels"),
			},
		},
	); err != nil {
		return
	}
	return
}

func (x *Service) SetTemplate(ctx context.Context, url string) (err error) {
	client := http.DefaultClient
	var req *http.Request
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = client.Do(req.WithContext(ctx)); err != nil {
		return
	}
	decoder := jsoniter.NewDecoder(resp.Body)
	for {
		var data map[string]interface{}
		if err = decoder.Decode(&data); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatalln(err)
			}
		}
	}
	return
}

func (x *Service) CodeKey(name string) string {
	return x.Values.KeyName("verify", name)
}

func (x *Service) CreateVerifyCode(ctx context.Context, name string, code string) error {
	return x.Redis.Set(ctx, x.CodeKey(name), code, time.Minute).Err()
}

// VerifyCode 校验验证码
func (x *Service) VerifyCode(ctx context.Context, name string, code string) (result bool, err error) {
	var exists int64
	if exists, err = x.Redis.Exists(ctx, x.CodeKey(name)).Result(); err != nil {
		return
	}
	if exists == 0 {
		return false, nil
	}
	var value string
	if value, err = x.Redis.Get(ctx, x.CodeKey(name)).Result(); err != nil {
		return
	}
	return value == code, nil
}

// RemoveVerifyCode 移除验证码
func (x *Service) RemoveVerifyCode(ctx context.Context, name string) error {
	return x.Redis.Del(ctx, x.CodeKey(name)).Err()
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
