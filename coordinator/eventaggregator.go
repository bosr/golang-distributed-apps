package coordinator

import "time"

// EventAggregator holds a map of registered functions for each listener
type EventAggregator struct {
	listeners map[string][]func(EventData)
}

// NewEventAggregator returns a new initialized EventAggregator
func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(EventData)),
	}
	return &ea
}

// AddListener adds a callback function
func (ea *EventAggregator) AddListener(name string, f func(EventData)) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

// PublishEvent calls all registered callback functions
func (ea *EventAggregator) PublishEvent(name string, eventData EventData) {
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
