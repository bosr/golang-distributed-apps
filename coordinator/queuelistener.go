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
	ea      *EventAggregator
}

// NewQueueListener returns a new QueueListener
func NewQueueListener(ea *EventAggregator) *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		ea:      ea,
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

// DiscoverSensors declares a new fanout exchange for broadcasting sensor name requests, and publishes an empty request, which is sufficient for them to reply.
func (ql *QueueListener) DiscoverSensors() {
	ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange, // name string,
		"fanout",                       // kind string,
		false,                          // durable bool,
		false,                          // autoDelete bool,
		false,                          // internal bool,
		false,                          // noWait bool,
		nil,                            // args amqp.Table,
	)

	ql.ch.Publish(
		qutils.SensorDiscoveryExchange, // exchange string,
		"",                // key string,
		false,             // mandatory bool,
		false,             // immediate bool,
		amqp.Publishing{}, // msg. amqp.Publishing,
	)
}

// ListenForNewSource does something right
func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch, true) // exchange will create the queue name for us
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

	ql.DiscoverSensors()

	fmt.Println("Listening for new sources.")
	for msg := range msgs {
		// new sensor has come online
		fmt.Println("New source discovered.")
		ql.ea.PublishEvent("DataSourceDiscovered", string(msg.Body))
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

		ed := EventData{
			Name:      sd.Name,
			Timestamp: sd.Timestamp,
			Value:     sd.Value,
		}

		ql.ea.PublishEvent("EventReceived_"+msg.RoutingKey, ed)
	}
}
