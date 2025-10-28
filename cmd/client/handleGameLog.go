package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func publishGameLog(ch *amqp.Channel, username, msg string) error {
	gamelog := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     msg,
		Username:    username,
	}

	err := pubsub.PublishGOB(
		ch,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.GameLogSlug, username),
		gamelog)
	if err != nil {
		return err
	}
	return nil
}
