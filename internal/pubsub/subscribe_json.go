package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	amqpChan, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		log.Fatalf("Queue does not exists or is not bound to the exchange: %s", err)
		return err
	}

	msgChan, err := amqpChan.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to create message channel: %s", err)
		return err
	}

	go func() {
		defer amqpChan.Close()
		for msg := range msgChan {
			var body_results T
			err := json.Unmarshal(msg.Body, &body_results)
			if err != nil {
				log.Fatalf("Error unmarshalling message from channel: %err", err)
			}
			handler(body_results)
		}
	}()

	return nil
}
