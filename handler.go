package main

import (
	"fmt"
	"reflect"
)

type internalHandlerAdapter interface {
	Execute(event any) error
}

type specificHandlerAdapter[T any] struct {
	// Stores the actual user function (e.g., handleOrderCreated)
	innerHandler func(event T) error
}

func (adapter *specificHandlerAdapter[T]) Execute(event any) error {
	if specificEvent, ok := event.(T); ok {
		return adapter.innerHandler(specificEvent)
	}
	return fmt.Errorf("type mismatch: expected %T but got %T", *new(T), event)
}

type handlerRegistry struct {
	handlers map[reflect.Type]internalHandlerAdapter
	store    EventStore
}

func newHandlerRegistry() *handlerRegistry {
	return &handlerRegistry{make(map[reflect.Type]internalHandlerAdapter), *NewEventStore()}
}

func Register[T any](registry *handlerRegistry, handler func(event T) error) {
	eventType := reflect.TypeOf(*new(T))
	adapter := &specificHandlerAdapter[T]{
		innerHandler: handler,
	}
	if registry.handlers == nil {
		registry.handlers = make(map[reflect.Type]internalHandlerAdapter)
	}
	registry.handlers[eventType] = adapter
	fmt.Printf("Registered handler for type %v\n", eventType)
}

func (r *handlerRegistry) Dispatch(event any) error {
	eventType := reflect.TypeOf(event)
	adapter, found := r.handlers[eventType]
	if !found {
		return fmt.Errorf("no handler for type %v", eventType)
	}
	return adapter.Execute(event)
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

func GetUserHandler(event Event) Event {
	getUserEvent, ok := event.(*GetUserEvent)
	if !ok {
		event.SetSuccess(false)
		return event
	}
	getUserEvent.User = pMap.getPlayer(getUserEvent.playerId)
	getUserEvent.SetSuccess(true)
	return getUserEvent
}

func GetSurroundingsHandler(event Event) Event {
	getSurroundingsEvent, ok := event.(*GetSurroundingsEvent)
	if !ok {
		event.SetSuccess(false)
		return event
	}
	getSurroundingsEvent.Minimap, _ = gMap.getSurroundingsOfPlayer(getSurroundingsEvent.OfPlayer()) //TODO: Handle bool flag
	getSurroundingsEvent.SetSuccess(true)
	return getSurroundingsEvent
}

func GetConfigHandler(event Event) Event {
	getConfigEvent, ok := event.(*GetConfigEvent)
	if !ok {
		event.SetSuccess(false)
		return event
	}
	getConfigEvent.Config = ConfigResponse{} //TODO: FIX ME
	getConfigEvent.SetSuccess(true)
	return getConfigEvent
}
