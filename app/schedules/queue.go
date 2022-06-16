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

func (x *Queue) Event(jobs *common.Jobs) (err error) {
	subj := x.Values.EventName("schedules")
	queue := x.Values.EventQueueName("schedules")
	var sub *nats.Subscription
	if _, err = x.Js.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		fmt.Printf(string(msg.Data))
		msg.Ack()
	}, nats.ManualAck()); err != nil {
		return
	}
	jobs.Store("schedules", sub)
	return
}
