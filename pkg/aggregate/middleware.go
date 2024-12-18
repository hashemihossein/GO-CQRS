package aggregate

import "context"

type Middleware interface {
	Execute(ctx context.Context, event AggregateEvent, next NextFunc) error
}

type NextFunc func(ctx context.Context, event AggregateEvent) error
