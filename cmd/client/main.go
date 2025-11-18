package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error loading .env file")
	}
	
	connectionString := os.Getenv("connectionString")
	
	conn,err := amqp.Dial(connectionString)
	if err != nil{
		log.Fatal("Error connecting to rabbitmq")
	}

	defer conn.Close()

	
	fmt.Println("Starting Peril client...")
}
