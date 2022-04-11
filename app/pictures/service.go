package pictures

import (
	"api/common"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Service struct {
	*common.Inject
}

func (x *Service) ImageInfo(ctx context.Context, url string) (result map[string]interface{}, err error) {
	var response *cos.Response
	if response, err = x.Cos.CI.Get(ctx, url, "imageInfo", nil); err != nil {
		return
	}
	if err = jsoniter.NewDecoder(response.Body).Decode(&result); err != nil {
		return
	}
	return
}
