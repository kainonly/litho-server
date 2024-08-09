package tencent

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/weplanx/go/help"
	"time"
)

func (x *Controller) CosPresigned(_ context.Context, c *app.RequestContext) {
	r, err := x.TencentX.CosPresigned()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, r)
}

func (x *Service) CosPresigned() (_ M, err error) {
	date := time.Now()
	expired := date.Add(time.Duration(x.V.TencentCosExpired) * time.Second)
	keyTime := fmt.Sprintf(`%d;%d`, date.Unix(), expired.Unix())
	name := help.Uuid()
	key := fmt.Sprintf(`%s/%s/%s`,
		x.V.Namespace, date.Format("20060102"), name)
	policy := M{
		"expiration": expired.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []interface{}{
			M{"bucket": x.V.TencentCosBucket},
			[]interface{}{"starts-with", "$key", key},
			M{"q-sign-algorithm": "sha1"},
			M{"q-ak": x.V.TencentSecretId},
			M{"q-sign-time": keyTime},
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
	return M{
		"key":              key,
		"policy":           policyText,
		"q-sign-algorithm": "sha1",
		"q-ak":             x.V.TencentSecretId,
		"q-key-time":       keyTime,
		"q-signature":      signature,
	}, nil
}
