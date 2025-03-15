package main

// Generic Event Handler
type EventHandler[T Event] interface {
	Handle(event T) error
}

// Generic Event Handler Function
type EventHandlerFunc[T Event] func(event T) error

// Handle calls the underlying function.
func (f EventHandlerFunc[T]) Handle(event T) error {
	return f(event)
}

// Example UserCreated Event Handler.
type UserCreatedHandler struct {
	// Dependencies (e.g., database connection)
}

// NewUserCreatedHandler creates a new UserCreatedHandler.
func NewUserCreatedHandler() *UserCreatedHandler {
	return &UserCreatedHandler{}
}

// Handle handles the UserCreated event.
func (h *UserCreatedHandler) Handle(event UserCreated) error {
	// Process the event (e.g., store user in database)
	// ...
	return nil
}

// Example of how to use the generic event handler function.
func HandleUserCreatedFunc(event UserCreated) error {
	//process the event
	return nil
}
