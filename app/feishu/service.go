package feishu

import (
	"api/common"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/weplanx/go/vars"
	"strings"
	"time"
)

type Service struct {
	*common.Inject
	Vars *vars.Service
}

func (x *Service) Decrypt(encrypt string, key string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", fmt.Errorf("base64StdEncode Error[%v]", err)
	}
	if len(buf) < aes.BlockSize {
		return "", errors.New("cipher  too short")
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
		return "", errors.New("ciphertext is not a multiple of the block size")
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
	key := x.Values.KeyName("feishu", "tenant_access_token")
	var exists int64
	if exists, err = x.Redis.Exists(ctx, key).Result(); err != nil {
		return
	}
	if exists == 0 {
		var id string
		if id, err = x.Vars.GetFeishuAppId(ctx); err != nil {
			return
		}
		var secret string
		if secret, err = x.Vars.GetFeishuAppSecret(ctx); err != nil {
			return
		}
		var body map[string]interface{}
		if _, err = x.HC.Feishu.R().
			SetBody(map[string]interface{}{
				"app_id":     id,
				"app_secret": secret,
			}).
			SetResult(&body).
			Post("/auth/v3/tenant_access_token/internal"); err != nil {
			return
		}
		if err = x.Redis.Set(ctx, key,
			body["tenant_access_token"],
			time.Second*time.Duration(body["expire"].(float64)),
		).Err(); err != nil {
			return
		}
	}
	return x.Redis.Get(ctx, key).Result()
}

type UserDto struct {
	AccessToken      string `json:"access_token"`
	AvatarUrl        string `json:"avatar_url"`
	AvatarThumb      string `json:"avatar_thumb"`
	AvatarMiddle     string `json:"avatar_middle"`
	AvatarBig        string `json:"avatar_big"`
	ExpiresIn        uint64 `json:"expires_in"`
	Name             string `json:"name"`
	EnName           string `json:"en_name"`
	OpenId           string `json:"open_id"`
	TenantKey        string `json:"tenant_key"`
	RefreshExpiresIn uint64 `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
}

func (x *Service) GetAccessToken(ctx context.Context, code string) (_ UserDto, err error) {
	var token string
	if token, err = x.GetTenantAccessToken(ctx); err != nil {
		return
	}
	var body struct {
		Code uint64  `json:"code"`
		Msg  string  `json:"msg"`
		Data UserDto `json:"data"`
	}
	if _, err = x.HC.Feishu.R().
		SetAuthToken(token).
		SetBody(map[string]interface{}{
			"grant_type": "authorization_code",
			"code":       code,
		}).
		SetResult(&body).
		Post("/authen/v1/access_token"); err != nil {
		return
	}
	if body.Code != 0 {
		err = errors.New(body.Msg)
		return
	}
	return body.Data, nil
}
