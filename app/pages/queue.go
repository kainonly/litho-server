package pages

import (
	"api/common"
	"fmt"
	"github.com/nats-io/nats.go"
)

type Queue struct {
	*common.Inject

	Service *Service
}

func (x *Queue) Event(subs *common.Subscriptions) (err error) {
	var sub *nats.Subscription
	if sub, err = x.Js.QueueSubscribe(
		x.Values.EventName("pages"),
		x.Values.KeyName("pages"),
		func(msg *nats.Msg) {
			fmt.Printf(string(msg.Data))
			msg.Ack()
		}, nats.ManualAck(),
	); err != nil {
		return
	}
	subs.Store("pages", sub)
	return
}
