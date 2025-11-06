package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any] (ch *amqp.Channel,exchange,key string,val T)error{
	jsonData, err := json.Marshal(val)
	if err != nil{
		fmt.Printf("error marshalling json: %v\n",err)
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body: jsonData,
	}
	err = ch.PublishWithContext(context.Background(),exchange,key,false,false,msg)
	return  err
}