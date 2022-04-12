package system

import (
	"api/common"
	"api/common/model"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/weplanx/go/helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) AppName() string {
	return x.Values.Namespace
}

type LoginLogDto struct {
	Time     time.Time          `bson:"time"`
	V        string             `bson:"v"`
	User     primitive.ObjectID `bson:"user"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	TokenId  string             `bson:"token_id"`
}

func NewLoginLogV10(data model.User, jti string) *LoginLogDto {
	return &LoginLogDto{
		Time:     time.Now(),
		V:        "v1.0",
		User:     data.ID,
		Username: data.Username,
		Email:    data.Email,
		TokenId:  jti,
	}
}

func (x *Service) WriteLoginLog(ctx context.Context, doc *LoginLogDto) (err error) {

	if _, err = x.Db.Collection("login_logs").InsertOne(ctx, doc); err != nil {
		return
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
