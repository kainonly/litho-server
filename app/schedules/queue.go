package schedules

import (
	"api/common"
	"fmt"
	"github.com/nats-io/nats.go"
)

type Queue struct {
	*common.Inject
	Js nats.JetStreamContext

	Service *Service
}

func (x *Queue) Event(subs *common.Subscriptions) (err error) {
	subj := x.Values.EventName("schedules")
	queue := x.Values.KeyName("schedules")
	var sub *nats.Subscription
	if sub, err = x.Js.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		fmt.Printf(string(msg.Data))
		msg.Ack()
	}, nats.ManualAck()); err != nil {
		return
	}
	subs.Store("schedules", sub)
	return
}
