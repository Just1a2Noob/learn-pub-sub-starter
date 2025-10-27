package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		log.Fatalf("Error marshalling value to json: %s", err)
		return err
	}

	publish_json := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, publish_json)
}
