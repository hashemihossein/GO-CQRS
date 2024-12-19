package command

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type CommandBus struct {
	handlers    map[reflect.Type]func(Command) error
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
var once sync.Once

func GetCommandBus() CommandBusInterface {
	once.Do(func() {
		commandBusInstance = &CommandBus{
			handlers:    make(map[reflect.Type]func(Command) error),
			middlewares: make([]Middleware, 0),
		}
	})

	return commandBusInstance
}

func (cb *CommandBus) RegisterCommandHandler(handler CommandHandler) error {
	cb.rwMu.Lock()
	defer cb.rwMu.Unlock()

	commandType := handler.GetCommandType()
	if _, exists := cb.handlers[commandType]; exists {
		return errors.New("this command handler has been registered before")
	}

	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	method, exists := handlerType.MethodByName("Handle")
	if !exists {
		return errors.New("handler does not have a Handle method")
	}

	if method.Type.NumIn() != 2 || method.Type.NumOut() != 1 {
		return errors.New("Handle method has incorrect signature")
	}

	specificCmdType := method.Type.In(1)
	if specificCmdType != commandType {
		return errors.New("Handle method parameter type does not match GetCommandType return type")
	}

	if method.Type.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.New("Handle method must return error")
	}

	wrapper := func(cmd Command) error {
		if reflect.TypeOf(cmd) != commandType {
			return errors.New("invalid command type")
		}

		results := handlerValue.MethodByName("Handle").Call([]reflect.Value{reflect.ValueOf(cmd)})

		if len(results) != 1 {
			return errors.New("Handle method should return exactly one value")
		}

		if errInterface := results[0].Interface(); errInterface != nil {
			return errInterface.(error)
		}

		return nil
	}

	cb.handlers[commandType] = wrapper
	return nil
}

func (cb *CommandBus) Dispatch(cmd Command) error {
	cb.rwMu.RLock()
	_, exists := cb.handlers[reflect.TypeOf(cmd)]
	cb.rwMu.RUnlock()
	if !exists {
		return errors.New("there is no registered handler for this command")
	}

	ctx := context.Background()
	return cb.DispatchWithContext(ctx, cmd)
}

func (cb *CommandBus) DispatchWithContext(ctx context.Context, cmd Command) error {
	cb.rwMu.RLock()
	handlerFunc, exists := cb.handlers[reflect.TypeOf(cmd)]
	middlewares := make([]Middleware, len(cb.middlewares))
	copy(middlewares, cb.middlewares)
	cb.rwMu.RUnlock()

	if !exists {
		return errors.New("there is no registered handler for this command")
	}

	chainFunc := func(ctx context.Context, cmd Command) error {
		return handlerFunc(cmd)
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		next := chainFunc
		chainFunc = func(ctx context.Context, cmd Command) error {
			return mw.Execute(ctx, cmd, next)
		}
	}

	return chainFunc(ctx, cmd)
}

func (cb *CommandBus) DispatchWithoutMiddlewares(cmd Command) error {
	cb.rwMu.RLock()
	handlerFunc, exists := cb.handlers[reflect.TypeOf(cmd)]
	cb.rwMu.RUnlock()
	if !exists {
		return errors.New("there is no registered handler for this command")
	}

	return handlerFunc(cmd)
}

func (cb *CommandBus) UseMiddleware(mw Middleware) {
	cb.rwMu.Lock()
	defer cb.rwMu.Unlock()
	cb.middlewares = append(cb.middlewares, mw)
}
