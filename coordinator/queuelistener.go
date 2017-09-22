package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/bosr/golang-distributed-apps/dto"
	"github.com/bosr/golang-distributed-apps/qutils"
	"github.com/streadway/amqp"
)

const (
	url = "amqp://guest:guest@localhost:5672"
)

// QueueListener does something good
type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
}

// NewQueueListener returns a new QueueListener
func NewQueueListener() *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

// ListenForNewSource does something right
func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch) // exchange will create the queue name for us
	// Re-bind the queue to the fanout exchange
	ql.ch.QueueBind(
		q.Name,       // name string
		"",           // key string, ignored by fanout exchange
		"amq.fanout", // exchange string,
		false,        // noWait bool, queue exists so cannot fail
		nil,          // args amqp.Table,
	)

	msgs, _ := ql.ch.Consume(
		q.Name, // queue string,
		"",     // consumer string,
		true,   // autoAck bool,
		false,  // exclusive bool,
		false,  // noLocal bool,
		false,  // noWait bool,
		nil,    // args amqp.Table,
	)

	for msg := range msgs {
		// new sensor has come online
		sourceChan, _ := ql.ch.Consume(
			string(msg.Body), // queue string,
			"",               // consumer string,
			true,             // autoAck bool,
			false,            // exclusive bool,
			false,            // noLocal bool,
			false,            // noWait bool,
			nil,              // args amqp.Table,
		)

		if ql.sources[string(msg.Body)] == nil {
			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

// AddListener does something good
func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("Received message: %v\n", sd)
	}
}
