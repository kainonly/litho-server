package app

import (
	"server/app/schedules"
	"server/common"
	"sync"
)

func SetJobs(
	schedules schedules.Queue,
) (jobs *common.Jobs, err error) {
	jobs = &sync.Map{}
	go schedules.Event(jobs)
	return
}
