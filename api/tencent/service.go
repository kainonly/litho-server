package tencent

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/utils"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/server/common"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	*common.Inject
}

// CosClient Cos 对象存储客户端
func (x *Service) CosClient(ctx context.Context) (client *cos.Client, err error) {
	var u *url.URL
	if u, err = url.Parse(
		fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`,
			x.V.TencentCosBucket, x.V.TencentCosRegion,
		),
	); err != nil {
		return
	}
	client = cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  x.V.TencentSecretId,
			SecretKey: x.V.TencentSecretKey,
		},
	})
	return
}

// CosPresigned 对象存储预签名
func (x *Service) CosPresigned(ctx context.Context) (result interface{}, err error) {
	date := time.Now()
	expired := date.Add(time.Duration(x.V.TencentCosExpired) * time.Second)
	keyTime := fmt.Sprintf(`%d;%d`,
		date.Unix(), expired.Unix())
	name, _ := gonanoid.Nanoid()
	key := fmt.Sprintf(`%s/%s/%s`,
		x.V.Namespace,
		date.Format("20060102"),
		name,
	)
	policy := map[string]interface{}{
		"expiration": expired.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []interface{}{
			map[string]interface{}{"bucket": x.V.TencentCosBucket},
			[]interface{}{"starts-with", "$key", key},
			map[string]interface{}{"q-sign-algorithm": "sha1"},
			map[string]interface{}{"q-ak": x.V.TencentSecretId},
			map[string]interface{}{"q-sign-time": keyTime},
		},
	}
	var policyText []byte
	if policyText, err = sonic.Marshal(policy); err != nil {
		return
	}
	signKeyHash := hmac.New(sha1.New, []byte(x.V.TencentSecretKey))
	signKeyHash.Write([]byte(keyTime))
	signKey := hex.EncodeToString(signKeyHash.Sum(nil))
	stringToSignHash := sha1.New()
	stringToSignHash.Write(policyText)
	stringToSign := hex.EncodeToString(stringToSignHash.Sum(nil))
	signatureHash := hmac.New(sha1.New, []byte(signKey))
	signatureHash.Write([]byte(stringToSign))
	signature := hex.EncodeToString(signatureHash.Sum(nil))
	return utils.H{
		"key":              key,
		"policy":           policyText,
		"q-sign-algorithm": "sha1",
		"q-ak":             x.V.TencentSecretId,
		"q-key-time":       keyTime,
		"q-signature":      signature,
	}, nil
}

// ImageInfo 图片信息
func (x *Service) ImageInfo(ctx context.Context, url string) (result map[string]interface{}, err error) {
	var client *cos.Client
	if client, err = x.CosClient(ctx); err != nil {
		return
	}
	var response *cos.Response
	if response, err = client.CI.Get(ctx, url, "imageInfo", nil); err != nil {
		return
	}
	b, _ := io.ReadAll(response.Body)
	if err = sonic.Unmarshal(b, &result); err != nil {
		return
	}
	return
}
