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

// Base event

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

// Additional events

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

type GetUserEvent struct {
	BaseEvent
	User Player
}

func (event GetUserEvent) Type() string {
	return "GetUserEvent"
}

func (event *GetUserEvent) SetSuccess(s bool) {
	event.success = s
}

type GetSurroundingsEvent struct {
	BaseEvent
	Minimap Surroundings
}

func (event GetSurroundingsEvent) Type() string {
	return "GetSurroundingsEvent"
}

func (event *GetSurroundingsEvent) SetSuccess(s bool) {
	event.success = s
}

type GetConfigEvent struct {
	BaseEvent
	Config ConfigResponse
}

func (event GetConfigEvent) Type() string {
	return "GetConfigEvent"
}

func (event *GetConfigEvent) SetSuccess(s bool) {
	event.success = s
}
