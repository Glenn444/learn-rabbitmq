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

	//declare and bind to a queue
	c,_,err := pubsub.DeclareAndBind(conn,"peril_topic","game_logs","game_logs.*",pubsub.Durable)

	//c, err := conn.Channel()
	if err != nil{
		fmt.Printf("error occurred connecting to amqp: %v", err)
	}
	// value := routing.PlayingState{
	// 	IsPaused: true,
	// }
	// pubsub.PublishJSON(c,routing.ExchangePerilDirect,routing.PauseKey,value)

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
				fmt.Printf("Hey, I'm pausing the game")
				value := routing.PlayingState{
					IsPaused: true,
				}
				pubsub.PublishJSON(c, routing.ExchangePerilDirect, routing.PauseKey, value)
			case "resume":
				fmt.Printf("Hey, I'm resuming the game")
				value := routing.PlayingState{
					IsPaused: false,
				}
				pubsub.PublishJSON(c,routing.ExchangePerilDirect,routing.PauseKey,value)
			case "quit":
				fmt.Printf("Hey, I'm exiting the game")
				return
			default:
				fmt.Printf("don't understand the command")
			}
		}
	}
}
