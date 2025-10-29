package pubsub

import (
	"bytes"
	"encoding/gob"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGOB[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
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

			buf := bytes.NewBuffer(msg.Body)
			var decoded_body T
			decoder := gob.NewDecoder(buf)

			if err := decoder.Decode(&decoded_body); err != nil {
				log.Fatalf("Error decoding message body: %s", err)
				msg.Nack(false, true)
			}

			acktype := handler(decoded_body)
			switch acktype {
			case Ack:
				msg.Ack(false)
			case NackRequeue:
				msg.Nack(false, true)
			case NackDiscard:
				msg.Nack(false, false)
			}
		}
	}()

	return nil
}
