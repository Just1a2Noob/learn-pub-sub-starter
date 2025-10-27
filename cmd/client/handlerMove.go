package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
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
			fmt.Print("Acknowledge\n")
			return pubsub.Ack
		default:
			fmt.Print("NackDiscard\n")
			return pubsub.NackDiscard
		}
	}
}
