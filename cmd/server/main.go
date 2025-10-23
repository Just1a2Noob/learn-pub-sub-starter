package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	// Using signal channels, we can listen for signals i.e. Ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Were waiting for a value to receive then active the code below it
	<-signalChan

	fmt.Printf("\nClossing connection...")
}
