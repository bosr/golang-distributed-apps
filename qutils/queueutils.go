package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// GetChannel returns a Connection and Channel to the broker
// at the provided url
func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection to message broker at %s", url)
	ch, err := conn.Channel()
	failOnError(err, "Failed to get channel for connection")

	return conn, ch
}

func getQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name string
		false, // durable bool
		false, // autoDelete bool
		false, // exclusive bool
		false, // noWait bool
		nil,   // args amqp.Table
	)

	failOnError(err, "Failed to declare queue")

	return &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
