package main

import (
	"fmt"
)

type handlerRegistry struct {
	handlers map[string]func(event Event) Event
}

func newHandlerRegistry() *handlerRegistry {
	return &handlerRegistry{make(map[string]func(event Event) Event)}
}

func (registry handlerRegistry) AddHandler(typeName string, handler func(event Event) Event) {
	registry.handlers[typeName] = handler
}

func (registry handlerRegistry) Handle(event Event) Event {
	return registry.handlers[event.Type()](event)
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

/*
func main() {
	registry := newHandlerRegistry()
	registry.AddHandler(CreateUserEvent{}.Type(), CreateUserHandler)
	baseEvent := BaseEvent{"playerId", time.Now(), BaseEvent{}.Type(), false}
	createUserEvent := CreateUserEvent{baseEvent, "username"}
	fmt.Println(baseEvent.Type(), createUserEvent.Type())
	registry.Handle(createUserEvent)
	} */
