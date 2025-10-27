package pubsub

import (
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

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

	isTransient := queueType == SimpleQueueTransient
	queueClassic := QueueClassicQuorum
	if isTransient {
		queueClassic = QueueClassicType
	}

	queue, err := connChan.QueueDeclare(
		queueName,
		queueType == SimpleQueueDurable,
		isTransient,
		isTransient,
		false,
		amqp.Table{
			"x-dead-letter-exchange": routing.ExchangePerilDeadLetter,
			"x-queue-type":           queueClassic,
		},
	)

	if err != nil {
		log.Fatalf("Error creating new queue to channel: %s", err)
	}

	err = connChan.QueueBind(queueName, key, exchange, false, nil)

	if err != nil {
		log.Fatalf("Error binding queue to exchange: %s", err)
	}

	return connChan, queue, nil
}
