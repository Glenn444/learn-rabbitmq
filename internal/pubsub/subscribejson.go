package pubsub

import (
	"encoding/json"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueType,
	handler func(T),
) error {
	ch, q, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		log.Fatalf("error occurred on DeclareAndBind %v\n", err)
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("messages.consume: %v", err)
	}

	go func() {
		for message := range messages {
			var data T
			json.Unmarshal(message.Body, &data)
			handler(data)

			message.Ack(false)

		}
	}()

	return nil
}
