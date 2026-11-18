package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	TransientType SimpleQueueType = "transient"
	DurableType   SimpleQueueType = "durable"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err := ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: data}); err != nil {
		return err
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return ch, amqp.Queue{}, err
	}
	durable := queueType == DurableType
	autoDelete := queueType == TransientType
	exclusive := queueType == TransientType
	q, err := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)
	if err != nil {
		return ch, amqp.Queue{}, err
	}
	if err := ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return ch, q, err
	}
	return ch, q, nil
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	ch, q, err := DeclareAndBind(conn, exchange, queueName, key, queueType) // return existing otherwise create.
	if err != nil {
		return err
	}
	deliveryCh, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for delivery := range deliveryCh {
			var data T
			if err := json.Unmarshal(delivery.Body, &data); err != nil {
				fmt.Println("Fail to unmarshal the q message.")
				delivery.Nack(false, false)
				continue
			}
			handler(data)
			delivery.Ack(false)
		}
	}()
	return nil
}
