package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {

}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
}

func failOnError(err Error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}
