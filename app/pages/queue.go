package pages

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
)

type Queue struct {
	Service *Service
}

func (x *Queue) Event(consumer pulsar.Consumer, msg pulsar.Message) {
	fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
		msg.ID(), string(msg.Payload()))
	consumer.Ack(msg)
}
