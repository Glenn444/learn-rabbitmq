package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := os.Getenv("amqpConnectionUrl")
	fmt.Printf("connection url: %v", connectionString)
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		log.Fatal("Error connecting to rabbitmq")
	}

	defer conn.Close()

	fmt.Println("Starting Peril client...")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error getting username: %v", err)
	}

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Fatalf("error binding to the queue: %v", err)
	}
	c, err := conn.Channel()
	if err != nil {
		fmt.Printf("error occurred connecting to amqp: %v", err)
	}
	// Keep the client running
	gamelogic.PrintServerHelp()

	for {
		input := gamelogic.GetInput()
		if input == nil {
			break
		}

		//Handle the commands
		if len(input) > 0 {
			switch input[0] {
			case "pause":
				fmt.Printf("sending a pause message")
				value := routing.PlayingState{
					IsPaused: true,
				}
				pubsub.PublishJSON(c, routing.ExchangePerilDirect, routing.PauseKey, value)
			case "resume":
				fmt.Printf("sending a resume message")
				value := routing.PlayingState{
					IsPaused: false,
				}
				pubsub.PublishJSON(c,routing.ExchangePerilDirect,routing.PauseKey,value)
			case "quit":
				fmt.Printf("exiting...")
				return
			default:
				fmt.Printf("don't understand the command")
			}
		}
	}
}
