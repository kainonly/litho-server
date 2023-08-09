package lark

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/go-resty/resty/v2"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
	Sessions *sessions.Service
	Locker   *locker.Locker
	Passport *passport.Passport
	Index    *index.Service
}

// Lark BaseURL https://open.larksuite.com/open-apis
var client = resty.New().
	SetTimeout(time.Second * 5).
	SetBaseURL("https://open.feishu.cn/open-apis")

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

func (x *Service) GetTenantAccessToken(ctx context.Context) (token string, err error) {
	key := x.V.Name("lark", "tenant_access_token")
	var exists int64
	if exists, err = x.RDb.Exists(ctx, key).Result(); err != nil {
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
			SetContext(ctx).
			SetBody(map[string]interface{}{
				"app_id":     x.V.LarkAppId,
				"app_secret": x.V.LarkAppSecret,
			}).
			SetResult(&result).
			Post("/auth/v3/tenant_access_token/internal"); err != nil {
			return
		}

		if err = x.RDb.Set(ctx, key,
			result.TenantAccessToken,
			time.Second*time.Duration(result.Expire)).Err(); err != nil {
			return
		}
	}
	return x.RDb.Get(ctx, key).Result()
}

func (x *Service) GetUserAccessToken(ctx context.Context, code string) (_ model.UserLark, err error) {
	var token string
	if token, err = x.GetTenantAccessToken(ctx); err != nil {
		return
	}
	var result struct {
		Code uint64         `json:"code"`
		Msg  string         `json:"msg"`
		Data model.UserLark `json:"data"`
	}
	if _, err = client.R().
		SetContext(ctx).
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

func (x *Service) Link(ctx context.Context, userId string, data model.UserLark) (_ *mongo.UpdateResult, err error) {
	id, _ := primitive.ObjectIDFromHex(userId)
	return x.Db.Collection("users").UpdateOne(ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"lark": data}},
	)
}

type LoginResult struct {
	User        model.User
	AccessToken string
}

func (x *Service) Login(ctx context.Context, openId string) (r *LoginResult, err error) {
	r = new(LoginResult)
	if r.User, err = x.Index.Logining(ctx, bson.M{"lark.open_id": openId, "status": true}); err != nil {
		return
	}

	if err = x.Db.Collection("users").
		FindOne(ctx, bson.M{"lark.open_id": openId, "status": true}).Decode(&r.User); err != nil {
		if err == mongo.ErrNoDocuments {
			err = common.ErrLoginNotExists
			return
		}
		return
	}

	userId := r.User.ID.Hex()
	if r.AccessToken, err = x.Index.CreateAccessToken(ctx, userId); err != nil {
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
		SetContext(ctx).
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
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(&result).
		Get("/task/v1/tasks"); err != nil {
		return
	}
	return
}
