package command

import "context"

type CommandMiddleware interface {
	Execute(ctx context.Context, cmd Command, next NextFunc) error
}

type NextFunc func(ctx context.Context, cmd Command) error
