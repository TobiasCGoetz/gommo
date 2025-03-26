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
	SetSuccess(bool)
}

type BaseEvent struct {
	playerId  string
	timestamp time.Time
	eventType string
	success   bool
}

func (event *BaseEvent) SetSuccess(s bool) {
	event.success = s
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

type CreateUserEvent struct {
	BaseEvent
	Username string
}

func (event CreateUserEvent) Type() string {
	return "CreateUserEvent"
}

func (event *CreateUserEvent) SetSuccess(s bool) {
	event.success = s
}
