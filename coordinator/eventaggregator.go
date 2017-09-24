package coordinator

import "time"

// EventRaiser slims down the listener interface
type EventRaiser interface {
	AddListener(eventName string, f func(interface{}))
}

// EventAggregator holds a map of registered functions for each listener
type EventAggregator struct {
	listeners map[string][]func(interface{})
}

// NewEventAggregator returns a new initialized EventAggregator
func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}
	return &ea
}

// AddListener adds a callback function
func (ea *EventAggregator) AddListener(name string, f func(interface{})) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

// PublishEvent calls all registered callback functions
func (ea *EventAggregator) PublishEvent(name string, eventData interface{}) {
	if ea.listeners[name] != nil {
		for _, r := range ea.listeners[name] {
			r(eventData) // call the callback with a copy not a pointer
		}
	}
}

// EventData is the local representation
type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}
