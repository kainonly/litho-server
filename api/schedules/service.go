package schedules

import (
	"github.com/weplanx/server/common"
	"github.com/weplanx/workflow"
	"sync"
)

type Service struct {
	*common.Inject
	M sync.Map
}

func (x *Service) Schedule(id string) (schedule *workflow.Schedule, err error) {
	if v, ok := x.M.Load(id); ok {
		schedule = v.(*workflow.Schedule)
		return
	}
	if schedule, err = x.Workflow.NewSchedule(id); err != nil {
		return
	}
	x.M.Store(id, schedule)
	return
}

func (x *Service) Ping(id string) (r bool, err error) {
	var schedule *workflow.Schedule
	if schedule, err = x.Schedule(id); err != nil {
		return
	}
	return schedule.Ping()
}

func (x *Service) Keys(id string) (keys []string, err error) {
	var schedule *workflow.Schedule
	if schedule, err = x.Schedule(id); err != nil {
		return
	}
	return schedule.Lists()
}
