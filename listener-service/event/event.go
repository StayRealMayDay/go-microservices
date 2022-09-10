package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declearExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         //durable
		false,        // aotudeleted
		false,        //internal
		false,        // no-wait
		nil,          //auguments
	)
}

func declearQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    //name
		false, // durable
		false, // delete when unsed
		true,  // exclusive
		false, // no wait
		nil,   //
	)
}
