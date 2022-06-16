package app

import (
	"api/app/schedules"
	"api/common"
	"sync"
)

func SetJobs(
	schedules schedules.Queue,
) (jobs *common.Jobs, err error) {
	jobs = &sync.Map{}
	go schedules.Event(jobs)
	return
}
