package app

import (
	"api/app/schedules"
	"api/common"
)

func Subscribe(
	schedules schedules.Queue,
) (subs *common.Subscriptions, err error) {
	subs = new(common.Subscriptions)
	go schedules.Event(subs)
	return
}
