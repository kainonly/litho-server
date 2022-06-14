package schedules

import (
	"api/common"
	"github.com/weplanx/schedule/client"
)

type Service struct {
	*common.Inject
	Client *client.Schedule
}
