package feishu

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/go-resty/resty/v2"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"github.com/weplanx/utils/locker"
	"github.com/weplanx/utils/passport"
	"github.com/weplanx/utils/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
	SessionsService *sessions.Service
	Locker          *locker.Locker
	Passport        *passport.Passport
}

var client = resty.New().
	SetBaseURL("https://open.feishu.cn/open-apis")

// Decrypt 解密
func (x *Service) Decrypt(encrypt string, key string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", errors.NewPublic("cipher  too short")
	}
	keyBs := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", fmt.Errorf("AESNewCipher Error[%v]", err)
	}
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.NewPublic("ciphertext is not a multiple of the block size")
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	n := strings.Index(string(buf), "{")
	if n == -1 {
		n = 0
	}
	m := strings.LastIndex(string(buf), "}")
	if m == -1 {
		m = len(buf) - 1
	}
	return string(buf[n : m+1]), nil
}

// GetTenantAccessToken 获取 TenantAccessToken
func (x *Service) GetTenantAccessToken(ctx context.Context) (token string, err error) {
	key := x.Values.Name("feishu", "tenant_access_token")
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var result struct {
			Code              uint64 `json:"code"`
			Msg               string `json:"msg"`
			TenantAccessToken string `json:"tenant_access_token"`
			Expire            int64  `json:"expire"`
		}
		if _, err = client.R().
			SetBody(map[string]interface{}{
				"app_id":     x.Values.FeishuAppId,
				"app_secret": x.Values.FeishuAppSecret,
			}).
			SetResult(&result).
			Post("/auth/v3/tenant_access_token/internal"); err != nil {
			return
		}
		if err = x.Redis.Set(ctx, key,
			result.TenantAccessToken,
			time.Second*time.Duration(result.Expire),
		).Err(); err != nil {
			return
		}
	}
	return x.Redis.Get(ctx, key).Result()
}

// GetUserAccessToken 获取 AccessToken
func (x *Service) GetUserAccessToken(ctx context.Context, code string) (_ model.FeishuUserData, err error) {
	var token string
	if token, err = x.GetTenantAccessToken(ctx); err != nil {
		return
	}
	var result struct {
		Code uint64               `json:"code"`
		Msg  string               `json:"msg"`
		Data model.FeishuUserData `json:"data"`
	}
	if _, err = client.R().
		SetAuthToken(token).
		SetBody(map[string]interface{}{
			"grant_type": "authorization_code",
			"code":       code,
		}).
		SetResult(&result).
		Post("/authen/v1/access_token"); err != nil {
		return
	}
	if result.Code != 0 {
		err = errors.NewPublic(result.Msg)
		return
	}
	return result.Data, nil
}

// Link 关联
func (x *Service) Link(ctx context.Context, userId string, data model.FeishuUserData) (_ *mongo.UpdateResult, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	return x.Db.Collection("users").UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"feishu": data}},
	)
}

// Login 免登陆
func (x *Service) Login(ctx context.Context, openId string, metadata *model.LoginMetadata) (ts string, err error) {
	var user model.User
	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{
			"feishu.open_id": openId,
			"status":         true,
		}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			err = errors.NewPublic("the user does not exist or has been frozen")
			return
		}

		return
	}

	userId := user.ID.Hex()
	metadata.Email = user.Email
	metadata.UserID = userId

	var maxLoginFailures bool
	if maxLoginFailures, err = x.Locker.Verify(ctx, userId, x.Values.LoginFailures); err != nil {
		return
	}
	if maxLoginFailures {
		err = errors.NewPublic("the user has exceeded the maximum number of login failures")
		return
	}

	jti, _ := gonanoid.Nanoid()
	metadata.TokenId = jti
	if ts, err = x.Passport.Create(userId, jti); err != nil {
		return
	}
	if err = x.Locker.Delete(ctx, userId); err != nil {
		return
	}
	if err = x.SessionsService.Set(ctx, userId, jti); err != nil {
		return
	}

	key := x.Values.Name("users", userId)
	if _, err = x.Redis.Del(ctx, key).Result(); err != nil {
		return
	}

	return
}

func (x *Service) CreateTask(ctx context.Context) (result map[string]interface{}, err error) {
	var token string
	if token, err = x.GetTenantAccessToken(ctx); err != nil {
		return
	}
	body := `{
		"summary": "每天喝八杯水，保持身心愉悦",
		"description": "多吃水果，多运动，健康生活，快乐工作。",
		"rich_summary": "富文本标题[飞书开放平台](https://open.feishu.cn)",
		"rich_description": "富文本备注[飞书开放平台](https://open.feishu.cn)",
		"extra": "dGVzdA==",
		"due": {
			"time": "1623124318",
			"timezone": "Asia/Shanghai",
			"is_all_day": false
		},
		"origin": {
			"platform_i18n_name": "{\"zh_cn\": \"IT 工作台\", \"en_us\": \"IT Workspace\"}",
			"href": {
				"url": "https://support.feishu.com/internal/foo-bar",
				"title": "反馈一个问题，需要协助排查"
			}
		},
		"can_edit":true,
		"custom": "{\"custom_complete\":{\"android\":{\"href\":\"https://www.feishu.cn/\",\"tip\":{\"zh_cn\":\"你好\",\"en_us\":\"hello\"}},\"ios\":{\"href\":\"https://www.feishu.cn/\",\"tip\":{\"zh_cn\":\"你好\",\"en_us\":\"hello\"}},\"pc\":{\"href\":\"https://www.feishu.cn/\",\"tip\":{\"zh_cn\":\"你好\",\"en_us\":\"hello\"}}}}",
		"follower_ids": ["ou_13585843f02bc94923ed17a007cbc9b1", "ou_219a0611de2a639aa939ee97013f37a5"],
		"collaborator_ids": ["ou_13585843f02bc94923ed17a007cbc9b1", "ou_219a0611de2a639aa939ee97013f37a5"],
		"repeat_rule": "FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR"
	}`

	if _, err = client.R().
		SetAuthToken(token).
		SetBody(body).
		SetResult(&result).
		Post("/task/v1/tasks"); err != nil {
		return
	}
	return
}

func (x *Service) GetTasks(ctx context.Context) (result map[string]interface{}, err error) {
	var token string
	if token, err = x.GetTenantAccessToken(ctx); err != nil {
		return
	}
	if _, err = client.R().
		SetAuthToken(token).
		SetResult(&result).
		Get("/task/v1/tasks"); err != nil {
		return
	}
	return
}
