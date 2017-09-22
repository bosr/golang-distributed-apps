package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/bosr/golang-distributed-apps/dto"
	"github.com/bosr/golang-distributed-apps/qutils"
	"github.com/streadway/amqp"
)

var (
	url = "amqp://guest:guest@localhost:5672"

	name     = flag.String("name", "sensor", "name of the sensor")
	freq     = flag.Uint("freq", 5, "update frequency in cycles/sec")
	max      = flag.Float64("max", 5., "max value for generated readings")
	min      = flag.Float64("min", 1., "min value for generated readings")
	stepSize = flag.Float64("step", 0.1, "max allowable change per sample")

	// rand gen seeded with nanosec
	r     = rand.New(rand.NewSource(time.Now().UnixNano()))
	value = r.Float64()*(*max-*min) + *min
	nom   = (*max-*min)/2 + *min
)

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch)

	// 5 cycles/sec -> 200 ms/cycle
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	fmt.Println(dur)
	signal := time.Tick(dur)
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		buf.Reset()
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"",             // exchange string
			dataQueue.Name, // key string
			false,          // mandatoory bool: because we are sure it exists
			false,          // immediate bool
			msg,            // msg amqp.Publishing
		)
		log.Printf("Reading sent. Value: %v\n", value)
	}
}

func calcValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}
