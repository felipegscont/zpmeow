package common

import (
	"context"
)

type Command interface {
	Validate() error
}

type CommandHandler[T Command] interface {
	Handle(ctx context.Context, cmd T) error
}

type Query interface {
	Validate() error
}

type QueryHandler[T Query, R any] interface {
	Handle(ctx context.Context, query T) (R, error)
}

type CommandBus interface {
	Execute(ctx context.Context, cmd Command) error
}

type QueryBus interface {
	Execute(ctx context.Context, query Query) (interface{}, error)
}
