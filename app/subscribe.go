package app

import (
	"api/app/pages"
	"api/common"
)

func Subscribe(
	pages pages.Queue,
) (subs *common.Subscriptions, err error) {
	subs = new(common.Subscriptions)
	go pages.Event(subs)
	return
}
