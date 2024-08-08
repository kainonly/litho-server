package lark

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/errors"
	"github.com/google/wire"
	"github.com/imroc/req/v3"
	"github.com/weplanx/go/locker"
	"github.com/weplanx/go/passport"
	"github.com/weplanx/go/sessions"
	"github.com/weplanx/server/api/index"
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/model"
	"strings"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	V        *common.Values
	Passport *passport.Passport

	LarkX  *Service
	IndexX *index.Service
}

type Service struct {
	*common.Inject
	Sessions *sessions.Service
	Locker   *locker.Locker
	Passport *passport.Passport
	IndexX   *index.Service
}

type M = map[string]interface{}

var client = req.C().
	SetBaseURL("https://open.feishu.cn/open-apis").
	SetJsonMarshal(sonic.Marshal).
	SetJsonUnmarshal(sonic.Unmarshal).
	SetTimeout(time.Second * 5)

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
		SetBearerAuthToken(token).
		SetBody(map[string]interface{}{
			"grant_type": "authorization_code",
			"code":       code,
		}).
		SetSuccessResult(&result).
		Post("/authen/v1/access_token"); err != nil {
		return
	}
	if result.Code != 0 {
		err = errors.NewPublic(result.Msg)
		return
	}
	return result.Data, nil
}

func (x *Service) GetTenantAccessToken(ctx context.Context) (token string, err error) {
	key := fmt.Sprintf(`%s:%s`, "lark", "tenant_access_token")
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
			SetBody(M{
				"app_id":     x.V.LarkAppId,
				"app_secret": x.V.LarkAppSecret,
			}).
			SetSuccessResult(&result).
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

type LoginResult struct {
	User        model.User
	AccessToken string
}
