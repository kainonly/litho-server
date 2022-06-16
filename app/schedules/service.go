package schedules

import (
	"api/common"
	"github.com/weplanx/schedule/client"
)

type Service struct {
	*common.Inject
	Client *client.Schedule
}

// GetKeys 获取调度服务已存在的标识
func (x *Service) GetKeys() (keys []string, err error) {
	return x.Client.List()
}
