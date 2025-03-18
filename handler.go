package main

import (
	"fmt"
	"time"
)

type handlerRegistry struct {
	handlers map[string]func(event BaseEvent)
}

func newHandlerRegistry() *handlerRegistry {
	return &handlerRegistry{make(map[string]func(event BaseEvent))}
}

func (registry handlerRegistry) AddHandler(typeName string, handler func())

func (registry handlerRegistry) Handle(event BaseEvent) {
	registry.handlers[event.EventType](event)
}

func ReadUserHandler(event ReadUser) {
	fmt.Println("Handling ReadUser event")
}

func main() {
	registry := newHandlerRegistry()
	registry.addHandler("ReadUser", ReadUserHandler)
	registry.Handle(ReadUser{"ReadUser0", time.Now(), "ReadUser", false})
}
