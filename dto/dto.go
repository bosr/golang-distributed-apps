package dto

import (
	"encoding/gob"
	"time"
)

// SensorMessage holds data to be transfered over the wire.
type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func init() {
	gob.Register(SensorMessage{})
}
