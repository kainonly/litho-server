package tencent

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/google/wire"
	tcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/weplanx/go/help"
	"github.com/weplanx/server/common"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var Provides = wire.NewSet(
	wire.Struct(new(Controller), "*"),
	wire.Struct(new(Service), "*"),
)

type Controller struct {
	TencentX *Service
}

type Service struct {
	*common.Inject
}

type M = map[string]interface{}

func (x *Service) Cos() (_ *cos.Client) {
	u, _ := url.Parse(fmt.Sprintf(`https://%s.cos.%s.myqcloud.com`, x.V.TencentCosBucket, x.V.TencentCosRegion))
	return cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  x.V.TencentSecretId,
			SecretKey: x.V.TencentSecretKey,
		},
	})
}

type TC3Option struct {
	Service   string
	Headers   map[string]string
	Timestamp int64
	Body      interface{}
}

func (x *Service) TC3Authorization(option TC3Option) string {
	algorithm := "TC3-HMAC-SHA256"
	canonicalURI := "/"
	canonicalQueryString := ""

	var keys []string
	for key := range option.Headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var canonicalHeaders string
	var signedHeaders string
	for _, key := range keys {
		k, v := strings.ToLower(key), strings.ToLower(option.Headers[key])
		canonicalHeaders += fmt.Sprintf("%s:%s\n", k, v)
		signedHeaders += ";" + k
	}

	signedHeaders = signedHeaders[1:]

	payload, _ := sonic.MarshalString(option.Body)
	hashedRequestPayload := common.Sha256hex(payload)
	canonicalRequest := fmt.Sprintf("POST\n%s\n%s\n%s\n%s\n%s",
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedRequestPayload,
	)

	date := time.Unix(option.Timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, option.Service)
	hashedCanonicalRequest := common.Sha256hex(canonicalRequest)
	string2sign := fmt.Sprintf("%s\n%d\n%s\n%s",
		algorithm,
		option.Timestamp,
		credentialScope,
		hashedCanonicalRequest,
	)

	secretDate := common.Hmacsha256(date, "TC3"+x.V.TencentSecretKey)
	secretService := common.Hmacsha256(option.Service, secretDate)
	secretSigning := common.Hmacsha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(common.Hmacsha256(string2sign, secretSigning)))

	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm,
		x.V.TencentSecretId,
		credentialScope,
		signedHeaders,
		signature,
	)
}

type KeyAuthResult struct {
	Date string
	Txt  string
}

func (x *Service) KeyAuth(source string, id string, key string) (r *KeyAuthResult, err error) {
	r = new(KeyAuthResult)
	location, _ := time.LoadLocation("Etc/UTC")
	r.Date = time.Now().In(location).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s\nx-source: %s", r.Date, source)

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	r.Txt = fmt.Sprintf("hmac id=\"%s\", algorithm=\"hmac-sha1\", headers=\"x-date x-source\", signature=\"%s\"",
		id, sign)
	return
}

type IpResult interface {
	GetMsg() string
	IsSuccess() bool
	GetDetail() interface{}
}

type Ipv4Result struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    struct {
		OrderNo string     `json:"orderNo"`
		Result  Ipv4Detail `json:"result"`
	} `json:"data"`
}

func (x *Ipv4Result) GetMsg() string {
	return x.Msg
}

func (x *Ipv4Result) IsSuccess() bool {
	return x.Success
}

func (x *Ipv4Result) GetDetail() interface{} {
	return x.Data.Result
}

type Ipv4Detail struct {
	Continent string `bson:"continent" json:"continent"`
	Country   string `bson:"country" json:"country"`
	Province  string `bson:"prov" json:"prov"`
	City      string `bson:"city" json:"city"`
	Owner     string `bson:"owner" json:"owner"`
	ISP       string `bson:"isp" json:"isp"`
	Areacode  string `bson:"areacode" json:"areacode"`
	Asnumber  string `bson:"asnumber" json:"asnumber"`
	Adcode    string `bson:"adcode" json:"adcode"`
	Zipcode   string `bson:"zipcode" json:"zipcode"`
	Timezone  string `bson:"timezone" json:"timezone"`
	Accuracy  string `bson:"accuracy" json:"accuracy"`
	Lat       string `bson:"lat" json:"lat"`
	Lng       string `bson:"lng" json:"lng"`
	Radius    string `bson:"radius" json:"radius"`
	Source    string `bson:"source" json:"source"`
}

func (x *Service) GetIpv4(ctx context.Context, ip string) (_ IpResult, err error) {
	source, kar := "market", new(KeyAuthResult)
	if kar, err = x.KeyAuth(source, x.V.IpSecretId, x.V.IpSecretKey); err != nil {
		return
	}
	var r *Ipv4Result
	var msg string
	if _, err = common.HttpClient(x.V.IpAddress).R().
		SetContext(ctx).
		SetHeader("X-Source", source).
		SetHeader("X-Date", kar.Date).
		SetHeader("Authorization", kar.Txt).
		SetQueryParam("ip", ip).
		SetSuccessResult(&r).
		SetErrorResult(&msg).
		Get("/ip/city/query"); err != nil {
		return nil, help.E("Tencent.GetIpv4Fail", err.Error())
	}
	if msg != "" {
		return nil, help.E("Tencent.GetIpv4Fail", msg)
	}
	return r, nil
}

type Ipv6Result struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    struct {
		OrderNo string     `json:"orderNo"`
		Result  Ipv6Detail `json:"result"`
	} `json:"data"`
}

func (x *Ipv6Result) GetMsg() string {
	return x.Msg
}

func (x *Ipv6Result) IsSuccess() bool {
	return x.Success
}

func (x *Ipv6Result) GetDetail() interface{} {
	return x.Data.Result
}

type Ipv6Detail struct {
	Continent string `bson:"continent" json:"continent"`
	Country   string `bson:"country" json:"country"`
	Province  string `bson:"prov" json:"province"`
	City      string `bson:"city" json:"city"`
	Owner     string `bson:"owner" json:"owner"`
	ISP       string `bson:"isp" json:"isp"`
	Areacode  string `bson:"areacode" json:"areacode"`
	Asnumber  string `bson:"asnumber" json:"asnumber"`
	Adcode    string `bson:"adcode" json:"adcode"`
	Zipcode   string `bson:"zipcode" json:"zipcode"`
	Timezone  string `bson:"timezone" json:"timezone"`
	Accuracy  string `bson:"accuracy" json:"accuracy"`
	Lat       string `bson:"lat" json:"lat"`
	Lng       string `bson:"lng" json:"lng"`
	Radius    string `bson:"radius" json:"radius"`
	Source    string `bson:"source" json:"source"`
}

func (x *Service) GetIpv6(ctx context.Context, ip string) (_ IpResult, err error) {
	source, kar := "market", new(KeyAuthResult)
	if kar, err = x.KeyAuth(source, x.V.Ipv6SecretId, x.V.Ipv6SecretKey); err != nil {
		return
	}

	var r *Ipv6Result
	var msg string
	if _, err = common.HttpClient(x.V.Ipv6Address).R().
		SetContext(ctx).
		SetHeader("X-Source", source).
		SetHeader("X-Date", kar.Date).
		SetHeader("Authorization", kar.Txt).
		SetQueryParam("ip", ip).
		SetSuccessResult(&r).
		SetErrorResult(&msg).
		Get("/ip/ipv6/query"); err != nil {
		return nil, help.E("Tencent.GetIpv6Fail", err.Error())
	}
	if msg != "" {
		return nil, help.E("Tencent.GetIpv6Fail", msg)
	}

	return r, nil
}

func (x *Service) SmsSend(ctx context.Context, sign string, tid string, params []string, phone []string) (err error) {
	credential := tcommon.NewCredential(
		x.V.SmsSecretId,
		x.V.SmsSecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, x.V.SmsRegion, cpf)
	request := sms.NewSendSmsRequest()
	request.SmsSdkAppId = tcommon.StringPtr(x.V.SmsAppId)
	request.SignName = tcommon.StringPtr(sign)
	request.TemplateId = tcommon.StringPtr(tid)
	request.TemplateParamSet = tcommon.StringPtrs(params)
	request.PhoneNumberSet = tcommon.StringPtrs(phone)
	request.SetContext(ctx)
	if _, err = client.SendSms(request); err != nil {
		return
	}
	return
}
