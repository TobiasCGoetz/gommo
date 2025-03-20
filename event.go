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
	return event.playerId
}

func (event BaseEvent) Type() string {
	return event.eventType
}

func (event BaseEvent) Timestamp() time.Time {
	return event.timestamp
}

type WriteUser struct {
	BaseEvent
	Username string
}

type ReadUser struct {
	BaseEvent
	Player Player
}
