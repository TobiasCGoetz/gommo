package main

import "time"

type BaseEvent struct {
	id        string
	timestamp time.Time
	eventType string
	success   bool
}

type WriteUser struct {
	BaseEvent
	Username string
}

type ReadUser struct {
	BaseEvent
	Player Player
}
