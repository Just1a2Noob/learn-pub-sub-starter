package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const Connection = "amqp://guest:guest@localhost:5672"

func main() {
	conn, err := amqp.Dial(Connection)
	if err != nil {
		log.Fatalf("Cannot connect to server %s", err)
		return
	}

	// To ensure the connection is closed when the program exits
	defer conn.Close()

	connChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error creating connection channel: %s", err)
	}

	fmt.Printf("Connection to server was successful")

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, "game_logs", "game_logs*", pubsub.SimpleQueueType{
		Durable:   true,
		Transient: false,
	})
	if err != nil {
		log.Fatalf("Error binding connection to queue: %s", err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(
				connChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resumes game state")
			err = pubsub.PublishJSON(
				connChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
