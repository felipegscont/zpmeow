package common

import (
	"context"
)

// Command represents a command in CQRS pattern
type Command interface {
	// Validate validates the command
	Validate() error
}

// CommandHandler represents a command handler
type CommandHandler[T Command] interface {
	// Handle processes the command
	Handle(ctx context.Context, cmd T) error
}

// Query represents a query in CQRS pattern
type Query interface {
	// Validate validates the query
	Validate() error
}

// QueryHandler represents a query handler
type QueryHandler[T Query, R any] interface {
	// Handle processes the query and returns the result
	Handle(ctx context.Context, query T) (R, error)
}

// CommandBus defines the contract for command dispatching
type CommandBus interface {
	// Execute executes a command
	Execute(ctx context.Context, cmd Command) error
}

// QueryBus defines the contract for query dispatching
type QueryBus interface {
	// Execute executes a query and returns the result
	Execute(ctx context.Context, query Query) (interface{}, error)
}
