package app

import (
	"github.com/weplanx/server/common"
	"github.com/weplanx/server/utils/errors"
	"time"
)

type Service struct {
	*common.Inject
}

func (x *Service) Index() time.Time {
	return time.Now()
}

func (x *Service) Test() error {
	return errors.NewPublic(1001, "错误啦")
}
