package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Were waiting for a value to receive then active the code below it
	<-signalChan

	fmt.Printf("\nClossing connection...")
}
