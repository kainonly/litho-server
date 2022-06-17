package schedules

import (
	"api/common"
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/nats-io/nats.go"
	"github.com/weplanx/go/engine"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
		var data engine.PublishDto
		if err := jsoniter.Unmarshal(msg.Data, &data); err != nil {
			log.Fatalln(err)
		}
		switch data.Event {
		case "create":
			msg.Ack()
			break
		case "update":
			oid, _ := primitive.ObjectIDFromHex(data.Id)
			if err := x.Service.Sync(context.TODO(), oid); err != nil {
				return
			}
			msg.Ack()
			break
		case "delete":
			if err := x.Service.Delete(data.Id); err != nil {
				return
			}
			msg.Ack()
			break
		}

	}, nats.ManualAck()); err != nil {
		return
	}
	jobs.Store("schedules", sub)
	return
}
