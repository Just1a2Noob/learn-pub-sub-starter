package pubsub

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType struct {
	Durable   bool
	Transient bool
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	connChan, err := conn.Channel()

	if err != nil {
		log.Fatalf("Error creating connection channel: %s", err)
	}

	queue, err := connChan.QueueDeclare(
		queueName,
		queueType.Durable,
		queueType.Transient,
		queueType.Transient,
		false,
		nil)

	if err != nil {
		log.Fatalf("Error creating new queue to channel: %s", err)
	}

	err = connChan.QueueBind(queueName, key, exchange, false, nil)

	if err != nil {
		log.Fatalf("Error binding queue to exchange: %s", err)
	}

	return connChan, queue, nil
}
