package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connStr := "amqp://guest:guest@localhost:5672/"
	brokerConnection, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer brokerConnection.Close()
	fmt.Println("RabbitMQ connection has been opened successfully...")
	rmqCh, err := brokerConnection.Channel()
	if err != nil {
		fmt.Println("Channel failed to open.")
		os.Exit(1)
	}
	if err := pubsub.PublishJSON(rmqCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}); err != nil {
		fmt.Printf("Failed to publish the data to the exchange. %s", err.Error())
		os.Exit(1)
	}
	fmt.Println("message has been published successfully.")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	_, ok := <-signalChan
	if ok {
		fmt.Println("Program is shutting down ...")
		os.Exit(0)
	}

}
