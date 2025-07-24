package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env file")
	}
}

func main() {
	StartConsumer()
	defer CloseConnections()

	msgs := ConsumeMessages()

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			parsedMessageEvent, _ := ParseMessageEvent(msg.Body)
			SendMessage(parsedMessageEvent)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	<-forever
}
