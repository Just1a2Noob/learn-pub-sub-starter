package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

const Connection = "amqp://guest:guest@localhost:5672"

func main() {
	conn, err := amqp.Dial(Connection)
	if err != nil {
		log.Fatalf("Cannot connect to server %s", err)
	}

	// To ensure the connection is closed when the program exits
	defer conn.Close()

	fmt.Printf("Connection to server was successful")

	// Using signal channels, we can listen for signals i.e. Ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Were waiting for a value to receive then active the code below it
	<-signalChan

	fmt.Printf("\nClossing connection...")
}
