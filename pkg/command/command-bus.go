package command

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type CommandBus struct {
	handlers    map[reflect.Type]CommandHandler
	middlewares []Middleware
	rwMu        sync.RWMutex
}

type CommandBusInterface interface {
	RegisterCommandHandler(handler CommandHandler) error
	Dispatch(cmd Command) error
	DispatchWithContext(ctx context.Context, cmd Command) error
	DispatchWithoutMiddlewares(cmd Command) error
	UseMiddleware(mw Middleware)
}

var commandBusInstance *CommandBus

func (cb *CommandBus) RegisterCommandHandler(handler CommandHandler) error {
	cb.rwMu.Lock()
	defer cb.rwMu.Unlock()

	commandType := handler.GetCommandName()
	if _, exists := cb.handlers[commandType]; exists {
		return errors.New("this handler have been registered before")
	}

	cb.handlers[commandType] = handler
	return nil
}

func (cb *CommandBus) Dispatch(cmd Command) error {
	cb.rwMu.RLock()
	defer cb.rwMu.RUnlock()

	commandType := reflect.TypeOf(cmd)
	_, exists := cb.handlers[commandType]
	if !exists {
		return errors.New("there is not registered handler for this command")
	}

	ctx := context.Background()
	return cb.DispatchWithContext(ctx, cmd)
}

func (cb *CommandBus) DispatchWithContext(ctx context.Context, cmd Command) error {
	commandType := reflect.TypeOf(cmd)
	handler := cb.handlers[commandType]

	chainFunc := func(ctx context.Context, cmd Command) error {
		return handler.Handle(cmd)
	}

	for i := len(cb.middlewares) - 1; i >= 0; i-- {
		mw := cb.middlewares[i]
		next := chainFunc
		chainFunc = func(ctx context.Context, cmd Command) error {
			return mw.Execute(ctx, cmd, next)
		}
	}

	return chainFunc(ctx, cmd)
}

func (cb *CommandBus) DispatchWithoutMiddlewares(cmd Command) error {
	cb.rwMu.RLock()
	defer cb.rwMu.RUnlock()

	commandType := reflect.TypeOf(cmd)
	handler, exists := cb.handlers[commandType]
	if !exists {
		return errors.New("there is not registered handler for this command")
	}

	return handler.Handle(cmd)
}

func (cb *CommandBus) UseMiddleware(mw Middleware) {
	cb.middlewares = append(cb.middlewares, mw)
}

var once sync.Once

func getCommandBus() CommandBusInterface {
	once.Do(func() {
		commandBusInstance = &CommandBus{
			handlers:    make(map[reflect.Type]CommandHandler),
			middlewares: make([]Middleware, 0),
		}
	})

	return commandBusInstance
}
