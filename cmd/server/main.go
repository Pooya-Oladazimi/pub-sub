package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"
	brokerConnection, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer brokerConnection.Close()
	fmt.Println("RabbitMQ connection has been opened successfully...")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	_, ok := <-signalChan
	if ok {
		fmt.Println("Program is shutting down ...")
		os.Exit(0)
	}
	fmt.Println("Starting Peril server...")
}
