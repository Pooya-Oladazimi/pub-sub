package main

import (
	"fmt"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Sending a pause message ...")
			if err := pubsub.PublishJSON(rmqCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true}); err != nil {
				fmt.Printf("Failed to publish the data to the exchange. %s", err.Error())
				os.Exit(1)
			}
			fmt.Println("Pause message has been published successfully.")
		case "resume":
			fmt.Println("Sending a resume message ...")
			if err := pubsub.PublishJSON(rmqCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false}); err != nil {
				fmt.Printf("Failed to publish the data to the exchange. %s", err.Error())
				os.Exit(1)
			}
			fmt.Println("Resume message has been published successfully.")
		case "quit":
			fmt.Println("Closing the server ...")
			os.Exit(0)
		default:
			fmt.Println("Please enter a valid command.")
			gamelogic.PrintServerHelp()
		}
	}

}
