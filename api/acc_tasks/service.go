package acc_tasks

import (
	"context"
	"github.com/weplanx/server/api/tencent"
	"github.com/weplanx/server/common"
	"net/url"
	"strconv"
	"time"
)

type Service struct {
	*common.Inject
	TencentService *tencent.Service
}

func (x *Service) Invoke(ctx context.Context) (r M, err error) {
	u, _ := url.Parse(x.V.AccelerateAddress)
	timestamp := time.Now().Unix()
	headerx := map[string]string{
		"Content-Type":  "application/json",
		"Host":          u.Host,
		"X-Scf-Cam-Uin": x.V.CamUin,
	}
	body := M{}
	if _, err = common.HttpClient(u.String()).R().SetContext(ctx).
		SetHeaders(headerx).
		SetHeader("X-Scf-Cam-Timestamp", strconv.FormatInt(timestamp, 10)).
		SetHeader("Authorization", x.TencentService.TC3Authorization(tencent.TC3Option{
			Service:   "scf",
			Headers:   headerx,
			Timestamp: timestamp,
			Body:      body,
		})).
		SetBody(body).
		SetSuccessResult(&r).
		Post(""); err != nil {
		return
	}

	return
}
