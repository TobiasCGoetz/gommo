package main

import (
	"fmt"
)

type handlerRegistry struct {
	handlers map[string]func(event Event) Event
	store    EventStore
}

func newHandlerRegistry() *handlerRegistry {
	return &handlerRegistry{make(map[string]func(event Event) Event), *NewEventStore()}
}

func (registry handlerRegistry) AddHandler(typeName string, handler func(event Event) Event) {
	registry.handlers[typeName] = handler
}

func (registry handlerRegistry) Handle(event Event) Event {
	registry.handlers[event.Type()](event)
	if event.Success() {
		registry.store.Append(event.ToJson())
	}
	return event
}

func CreateUserHandler(event Event) Event {
	createUserEvent, ok := event.(*CreateUserEvent)
	if !ok {
		event.SetSuccess(false)
		return event
	}
	event.SetSuccess(true)
	fmt.Println("Successfully handled ", createUserEvent.Type())
	return event
}
