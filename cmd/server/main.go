package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error Loading .env file")
	}

	connectionString := os.Getenv("amqpConnectionUrl")
	conn, err := amqp.Dial(connectionString)

	if err != nil {
		fmt.Printf("error occurred connecting to amqp: %v", err)
	}
	//fmt.Printf("connection was successful\n")
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil{
		fmt.Printf("error occurred connecting to amqp: %v", err)
	}
	value := routing.PlayingState{
		IsPaused: true,
	}
	pubsub.PublishJSON(c,routing.ExchangePerilDirect,routing.PauseKey,value)

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan

	fmt.Printf("signal to exit received: %v\n",sig)
}
