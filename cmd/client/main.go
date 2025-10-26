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

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error in getting username: %s", err)
	}

	// amqpChan, queue,
	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		strings.Join([]string{routing.PauseKey, username}, "."),
		routing.PauseKey,
		pubsub.SimpleQueueType{Transient: true})

	if err != nil {
		log.Fatalf("Error declaring and binding connection: %s", err)
	}

	gamestate := gamelogic.NewGameState(username)
	queue_name := fmt.Sprintf("pause.%s", username)

	err = pubsub.SubscribeJSON(
		conn,
		string(routing.ExchangePerilDirect),
		queue_name,
		string(routing.PauseKey),
		pubsub.SimpleQueueType{Transient: true},
		handlerPause(gamestate),
	)

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := gamestate.CommandSpawn(words)
			if err != nil {
				log.Printf("Error spawning units: %s", err)
			}
		case "move":
			_, err := gamestate.CommandMove(words)
			if err != nil {
				log.Printf("Error moving unit: %s", err)
			}
		case "status":
			gamestate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming is not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}
