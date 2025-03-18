package main

import "time"

type Event interface {
	ID() string
	Timestamp() time.Time
	Type() string
	Succeded() bool
}

type BaseEvent struct {
	id        string
	timestamp time.Time
	eventType string
	success   bool
}

// ID returns the event ID.
func (e BaseEvent) ID() string {
	return e.id
}

// Timestamp returns the event timestamp.
func (e BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

// Type returns the event type.
func (e BaseEvent) Type() string {
	return e.eventType
}

// Succeeded returns the success indicator
func (e BaseEvent) Succeeded() bool {
	return e.success
}

type UserCreated struct {
	BaseEvent
	Username string
}
