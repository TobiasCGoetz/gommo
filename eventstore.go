package main

import "fmt"

type EventStore struct {
	events [][]byte
}

func NewEventStore() *EventStore {
	return &EventStore{make([][]byte, 0)}
}

func (es *EventStore) Append(event []byte) {
	es.events = append(es.events, event)
}

func (es EventStore) Print() {
	for _, b := range es.events {
		fmt.Println("%t", b)
	}
}
