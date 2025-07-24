package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	exchange string = "megebase.topic"
	err      error
	conn     *amqp.Connection
	ch       *amqp.Channel
	q        amqp.Queue
)

type MessageEvent struct {
	Service string                 `json:"service"`
	Entity  string                 `json:"entity"`
	Action  string                 `json:"action"`
	Channel string                 `json:"channel,omitempty"`
	To      string                 `json:"to,omitempty"`
	Data    map[string]interface{} `json:"data"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func StartConsumer() {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Panic("RABBITMQ_URL not set")
	}

	conn, err = amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err = ch.QueueDeclare(
		"email", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

func ConsumeMessages() <-chan amqp.Delivery {
	q, err = ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,    // queue name
		"#.email", // routing key
		exchange,  // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}

func CloseConnections() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
