package tencent

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/bytedance/sonic/decoder"
	"github.com/weplanx/server/common"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	*common.Inject
}

type KeyAuthResult struct {
	Date string
	Txt  string
}

func (x *Service) KeyAuth(source string) (r *KeyAuthResult, err error) {
	r = new(KeyAuthResult)
	location, _ := time.LoadLocation("Etc/GMT")
	r.Date = time.Now().In(location).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s\nx-source: %s", r.Date, source)

	mac := hmac.New(sha1.New, []byte(x.V.IpSecretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	r.Txt = fmt.Sprintf("hmac id=\"%s\", algorithm=\"hmac-sha1\", headers=\"x-date x-source\", signature=\"%s\"",
		x.V.IpSecretId, sign)

	return
}

type CityResult struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Data    struct {
		OrderNo string `json:"orderNo"`
		Result  Detail `json:"result"`
	} `json:"data"`
}

func (x *CityResult) GetDetail() Detail {
	return x.Data.Result
}

type Detail struct {
	Continent string `bson:"continent" json:"continent"`
	Country   string `bson:"country" json:"country"`
	Prov      string `bson:"prov" json:"prov"`
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

func (x *Service) GetCity(ip string) (r *CityResult, err error) {
	source, kar := "market", new(KeyAuthResult)
	if kar, err = x.KeyAuth(source); err != nil {
		return
	}

	baseUrl, _ := url.Parse(x.V.IpAddress)
	u := baseUrl.JoinPath("/ip/city/query")
	query := u.Query()
	query.Add("ip", ip)
	query.Encode()
	u.RawQuery = query.Encode()

	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	req.Header.Set("X-Source", source)
	req.Header.Set("X-Date", kar.Date)
	req.Header.Set("Authorization", kar.Txt)

	client := &http.Client{Timeout: time.Second * 5}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}
	if err = decoder.NewStreamDecoder(res.Body).Decode(&r); err != nil {
		return
	}

	return
}
