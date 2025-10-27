package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, channel *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print(">")
		outcome := gs.HandleMove(move)

		switch outcome {
		case gamelogic.MoveOutcomeSamePlayer:
			fmt.Print("NackDiscard\n")
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			fmt.Print("Acknowledge\n")
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			warKey := fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, gs.GetUsername())
			err := pubsub.PublishJSON(
				channel,
				routing.ExchangePerilTopic,
				warKey,
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Printf("Failed to publish war message: %v\n", err)
				return pubsub.NackRequeue
			}

			fmt.Print("NackRequeue\n")
			return pubsub.NackDiscard
		default:
			fmt.Print("NackDiscard\n")
			return pubsub.NackDiscard
		}
	}
}
