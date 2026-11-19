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
	if err := pubsub.SubscribeJSON(brokerConnection, routing.ExchangePerilDirect, qName, routing.PauseKey, pubsub.TransientType, handlerPause(gameState)); err != nil {
		fmt.Println("Subscribe to pause queue failed.")
		os.Exit(1)
	}
	armyMoveRoutingKey := fmt.Sprintf("%s.*", routing.ArmyMovesPrefix)
	armyMoveQName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	if err := pubsub.SubscribeJSON(brokerConnection, routing.ExchangePerilTopic, armyMoveQName, armyMoveRoutingKey, pubsub.TransientType, handlerMove(gameState)); err != nil {
		fmt.Println("Subscribe to pause queue failed.")
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
			}
		case "move": // move an existing unit on the map
			mv, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err.Error())
			}
			ch, err := brokerConnection.Channel()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if err := pubsub.PublishJSON(ch, routing.ExchangePerilTopic, fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username), mv); err != nil {
				fmt.Println("Publishing move to exchange failed.")
				os.Exit(1)
			}
			fmt.Println("Move has been successfully published.")
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

func handlerMove(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		fmt.Println("Move has been detected.")
	}
}
