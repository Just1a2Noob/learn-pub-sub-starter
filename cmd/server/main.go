package main

import (
	"fmt"
	"log"
	"strings"

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

	gamelogic.PrintServerHelp()

	for {
		user_input := gamelogic.GetInput()

		fmt.Println(strings.ToLower(user_input[0]))

		if user_input[0] == "pause" {

			// Below publishes an exchange message
			err = pubsub.PublishJSON(
				connChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)

			if err != nil {
				log.Fatalf("Error publishing message: %s", err)
				return
			}
		}

		if user_input[0] == "resume" {
			// Below publishes an exchange message
			err = pubsub.PublishJSON(
				connChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)

			if err != nil {
				log.Fatalf("Error publishing message: %s", err)
				return
			}
		}

		if user_input[0] == "quit" {
			fmt.Printf("\nClossing connection...")
			break
		} else {
			fmt.Printf("\nCommand not understood\n")
		}
	}
}
