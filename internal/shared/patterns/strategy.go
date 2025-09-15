package patterns

import (
	"context"
)

type Strategy[T, R any] interface {
	Execute(ctx context.Context, input T) (R, error)
}

type StrategyFactory[T, R any] interface {
	CreateStrategy(strategyType string) Strategy[T, R]
}

type ProcessingStrategy[T any] interface {
	Validate(data T) error
	Process(ctx context.Context, data T) (T, error)
}

type SendingStrategy[T, R any] interface {
	Send(ctx context.Context, data T) (R, error)
}
