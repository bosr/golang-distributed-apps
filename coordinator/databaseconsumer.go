package coordinator

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/bosr/golang-distributed-apps/dto"
	"github.com/bosr/golang-distributed-apps/qutils"
	"github.com/streadway/amqp"
)

const (
	maxRate = 5
)

// DatabaseConsumer consumes events and sends them to the database after processing.
type DatabaseConsumer struct {
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	queue   *amqp.Queue
	sources []string
}

// NewDatabaseConsumer returns a new initialized DatabaseConsumer
func NewDatabaseConsumer(er EventRaiser) *DatabaseConsumer {
	dc := DatabaseConsumer{
		er: er,
	}

	dc.conn, dc.ch = qutils.GetChannel(url)
	dc.queue = qutils.GetQueue(qutils.PersistReadingsQueue, dc.ch, false)
	dc.er.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		dc.SubscribeToDataEvent(eventData.(string))
	})

	return &dc
}

// SubscribeToDataEvent subscribes the DatabaseConsumer to the events
func (dc *DatabaseConsumer) SubscribeToDataEvent(eventName string) {
	for _, v := range dc.sources {
		if v == eventName {
			return // we already subscribed
		}
	}

	dc.er.AddListener("MessageReceived_"+eventName, func() func(interface{}) {
		prevTime := time.Unix(0, 0)

		buf := new(bytes.Buffer)

		return func(eventData interface{}) {
			ed := eventData.(EventData)
			if time.Since(prevTime) > maxRate {
				prevTime = time.Now()

				sm := dto.SensorMessage{
					Name:      ed.Name,
					Value:     ed.Value,
					Timestamp: ed.Timestamp,
				}

				buf.Reset()

				enc := gob.NewEncoder(buf)
				enc.Encode(sm)

				msg := amqp.Publishing{
					Body: buf.Bytes(),
				}

				dc.ch.Publish(
					"", // exchange string,
					qutils.PersistReadingsQueue, // key string,
					false, // mandatory bool,
					false, // immediate bool,
					msg,   // msg amqp.Publishing
				)
			}
		}
	}())
}
