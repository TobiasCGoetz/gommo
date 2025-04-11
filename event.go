package main

import (
	"encoding/json"
	"time"
)

type Event interface {
	Timestamp() time.Time
	OfPlayer() string
	ToJson() []byte
	Success() bool
	SetSuccess(bool) Event
}

// Base event
type BaseEvent struct {
	playerId  string
	timestamp time.Time
	success   bool
}

func NewBaseEvent(playerId string) BaseEvent {
	return BaseEvent{playerId, time.Now(), false}
}

func (event BaseEvent) SetSuccess(s bool) Event {
	event.success = s
	return event
}

func (event BaseEvent) Success() bool {
	return event.success
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

func (event BaseEvent) ToJson() []byte {
	jsonData, err := json.Marshal(event)
	if err != nil {
		r, _ := json.Marshal("")
		return r
	}
	return jsonData
}

// Additional events

type CreateUserEvent struct {
	BaseEvent
	Username string
}

type GetUserEvent struct {
	BaseEvent
	User Player
}

type GetSurroundingsEvent struct {
	BaseEvent
	Minimap Surroundings
}

func NewGetSurroundingsEvent(playerId string) GetSurroundingsEvent {
	return GetSurroundingsEvent{NewBaseEvent(playerId), Surroundings{}}
}

type GetConfigEvent struct {
	BaseEvent
	Config ConfigResponse
}
