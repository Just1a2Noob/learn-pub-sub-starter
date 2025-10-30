package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
	unmarshaller func([]byte) (T, error),
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

			val, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("Failed to decode message: %s", err)
				continue
			}

			acktype := handler(val)
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

func SubscribeGOB[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
) error {
	return Subscribe(conn, exchange, queueName, key, queueType, handler, func(data []byte) (T, error) {
		var val T
		decoder := gob.NewDecoder(bytes.NewReader(data))
		err := decoder.Decode(&val)
		return val, err
	})
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) Acktype,
) error {
	return Subscribe(conn, exchange, queueName, key, queueType, handler, func(data []byte) (T, error) {
		var val T
		decoder := json.NewDecoder(bytes.NewBuffer(data))
		err := decoder.Decode(&val)
		return val, err
	})
}
