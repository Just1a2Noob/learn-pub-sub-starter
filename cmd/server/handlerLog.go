package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HandlerLog(ch *amqp.Channel) func(routing.GameLog) pubsub.Acktype {
	return func(logs routing.GameLog) pubsub.Acktype {
		defer fmt.Print("> ")

		err := gamelogic.WriteLog(logs)
		if err != nil {
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
