package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(recog gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Print("> ")

		outcome, winner, loser := gs.HandleWar(recog)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon, gamelogic.WarOutcomeYouWon, gamelogic.WarOutcomeDraw:
			var msg string
			if outcome == gamelogic.WarOutcomeDraw {
				msg = fmt.Sprintf("%s and %s resulted in a draw", winner, loser)
			} else {
				msg = fmt.Sprintf("%s won a war against %s", winner, loser)
			}

			err := publishGameLog(ch, gs.Player.Username, msg)
			if err != nil {
				return pubsub.NackRequeue
			}

			return pubsub.Ack
		default:
			fmt.Print("Error outcome of war cannot be determined")
			return pubsub.NackDiscard
		}
	}
}
