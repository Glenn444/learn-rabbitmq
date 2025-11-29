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
	
	gameState := gamelogic.NewGameState(username)
	// Keep the client running
	pubsub.SubscribeJSON(conn,routing.ExchangePerilDirect,username,routing.PauseKey,pubsub.Transient,handlerPause(gameState))
	for {
		input := gamelogic.GetInput()
		if input == nil {
			break
		}

		//Handle the commands
		if len(input) > 0 {
			switch input[0] {
			case "spawn":
				err := gameState.CommandSpawn(input)
				if err != nil{
					fmt.Printf("error occured: %v\n",err)
				}
			case "move":
				armyMove,err := gameState.CommandMove(input)
				if err != nil{
					log.Fatalf("error occured: %v\n",err)
				}
				fmt.Printf("move %v\n",armyMove.ToLocation)
			case "status":
				gameState.CommandStatus()
			case "help":
				gamelogic.PrintClientHelp()
			case "spam":
				fmt.Printf("Spamming not allowed yet\n")
			case "quit":
				gamelogic.PrintQuit()
				return
			default:
				fmt.Printf("Entered wrong command\n")
			}
		}
	}
}
