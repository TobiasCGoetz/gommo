package main

import (
	"encoding/json"
	"time"
)

type Event interface {
	Type() string
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
	eventType string
	success   bool
}

func NewBaseEvent(playerId string, eType string) BaseEvent {
	return BaseEvent{playerId, time.Now(), eType, false}
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

func (event CreateUserEvent) Type() string {
	return "CreateUserEvent"
}

type GetUserEvent struct {
	BaseEvent
	User Player
}

func (event GetUserEvent) Type() string {
	return "GetUserEvent"
}

type GetSurroundingsEvent struct {
	BaseEvent
	Minimap Surroundings
}

func NewGetSurroundingsEvent(playerId string) GetSurroundingsEvent {
	return GetSurroundingsEvent{NewBaseEvent(playerId, GetSurroundingsEvent{}.Type()), Surroundings{}}
}

func (event GetSurroundingsEvent) Type() string {
	return "GetSurroundingsEvent"
}

type GetConfigEvent struct {
	BaseEvent
	Config ConfigResponse
}

func (event GetConfigEvent) Type() string {
	return "GetConfigEvent"
}
