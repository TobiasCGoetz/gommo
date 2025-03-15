package main

import "time"

type Event interface {
	ID() string
	Timestamp() time.Time
	Payload() any
	Type() string
}

type BaseEvent struct {
	id        string
	timestamp time.Time
	payload   any
	eventType string
}

// ID returns the event ID.
func (e BaseEvent) ID() string {
	return e.id
}

// Timestamp returns the event timestamp.
func (e BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

// Payload returns the event payload.
func (e BaseEvent) Payload() interface{} {
	return e.payload
}

// Type returns the event type.
func (e BaseEvent) Type() string {
	return e.eventType
}

// Example event struct
type UserCreated struct {
	BaseEvent
	Username string
}
