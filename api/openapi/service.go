package openapi

import (
	"github.com/weplanx/openapi/client"
	"github.com/weplanx/server/common"
)

type Service struct {
	*common.Inject
}

func (x *Service) Client() (*client.Client, error) {
	return client.New(
		x.Values.OpenapiUrl,
		client.SetApiGateway(x.Values.OpenapiKey, x.Values.OpenapiSecret),
	)
}
