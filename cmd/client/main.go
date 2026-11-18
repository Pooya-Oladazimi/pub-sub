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
	fmt.Println("Starting Peril client...")
	connStr := "amqp://guest:guest@localhost:5672/"
	brokerConnection, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println("RabbitMq connection failed.")
		os.Exit(1)
	}
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	qName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	gameState := gamelogic.NewGameState(username)
	err = pubsub.SubscribeJSON(brokerConnection, routing.ExchangePerilDirect, qName, routing.PauseKey, pubsub.TransientType, handlerPause(gameState))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn": // add a new unit to a location on map
			err := gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err.Error())
				// os.Exit(1)
			}
		case "move": // move an existing unit on the map
			_, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err.Error())
				// os.Exit(1)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			os.Exit(0)
		default:
			fmt.Println("Unknown command.")
			continue
		}

	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// _, ok := <-signalChan
	// if ok {
	// fmt.Println("Program is shutting down ...")
	// os.Exit(0)
	// }

}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}
