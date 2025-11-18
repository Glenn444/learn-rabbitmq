package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {
		fmt.Printf("error marshalling json: %v\n", err)
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonData,
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
	return err
}

type SimpleQueType int

const (
	Durable SimpleQueType = iota
	Transient
)

func (q SimpleQueType) String() string {
	switch q {
	case Durable:
		return "durable"
	case Transient:
		return "transient"
	default:
		return "unknown"
	}
}

func (q SimpleQueType) IsDurable() bool {
	return q == Durable
}

func (q SimpleQueType) IsAutoDelete() bool {
	return q == Transient
}

func (q SimpleQueType) IsExclusive() bool {
	return q == Transient
}

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueType) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		queueType.IsDurable(),
		queueType.IsAutoDelete(),
		queueType.IsExclusive(),
		false,
		nil,
	)
	if err != nil{
		return nil, amqp.Queue{},err
	}

	err = ch.QueueBind(queueName,key,exchange,false,nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	return ch, q, nil

}
