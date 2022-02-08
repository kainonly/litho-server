package app

import (
	"api/app/pages"
	"api/common"
	"github.com/apache/pulsar-client-go/pulsar"
)

func Consumer(
	values *common.Values,
	client pulsar.Client,
	pages pages.Queue,
) (consumers map[string]pulsar.Consumer, err error) {
	if consumers["event"], err = client.Subscribe(pulsar.ConsumerOptions{
		Topic:            values.Pulsar.Topics["event"],
		SubscriptionName: "event",
		Type:             pulsar.Shared,
	}); err != nil {
		return
	}
	go func() {
		for c := range consumers["event"].Chan() {
			msg := c.Message
			switch msg.Key() {
			case "pages":
				pages.Event(consumers["event"], msg)
				break
			}
		}
	}()
	return
}
