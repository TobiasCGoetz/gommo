package main

import "time"

type Event interface {
	Type() string
	Timestamp() time.Time
	OfPlayer() string
}

type BaseEvent struct {
	playerId  string
	timestamp time.Time
	eventType string
	Success   bool
}

func (event BaseEvent) OfPlayer() string {
	return "BaseEvent"
}

func (event BaseEvent) Type() string {
	return "BaseEvent"
}

func (event BaseEvent) Timestamp() time.Time {
	return event.timestamp
}

type CreateUserEvent struct {
	BaseEvent
	Username string
}

func (event CreateUserEvent) Type() string {
	return "CreateUserEvent"
}
