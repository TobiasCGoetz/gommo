package main

import "time"

type BaseEvent struct {
	Id        string
	Timestamp time.Time
	EventType string
	Success   bool
}

type WriteUser struct {
	BaseEvent
	Username string
}

type ReadUser struct {
	BaseEvent
	Player Player
}
