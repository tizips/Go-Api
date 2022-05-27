package amqp

import (
	"github.com/gookit/goutil/dump"
	amqpPkg "github.com/streadway/amqp"
	"saas/kernel/queue/amqp"
)

var Amqp *amqp.Connection

var consumers *amqp.Connection

func Init() {
	go producer()
	go consumer()
}

func producer() {

	Amqp = new(amqp.Connection)

	Amqp.InitProducer(false)
}

func consumer() {

	consumers = new(amqp.Connection)

	consumers.Consumers = []amqp.Consumer{
		{
			Exchange: "demo",
			Queue:    "demo",
			Name:     "demo",
			Numbers:  3,
		},
	}

	consumers.ConsumeHandle = func(message amqpPkg.Delivery) {
		dump.P(string(message.Body))

		_ = message.Ack(true)
	}

	consumers.InitConsumer(false)
}
