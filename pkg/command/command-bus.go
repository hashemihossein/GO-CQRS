package command

import (
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
	Execute(cmd Command) error
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

func (cb *CommandBus) Execute(cmd Command) error {
	cb.rwMu.RLock()
	defer cb.rwMu.RUnlock()

	commandType := reflect.TypeOf(cmd)
	handler, exists := cb.handlers[commandType]
	if !exists {
		return errors.New("there is not registered handler for this command")
	}

	return handler.Handle(cmd)
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
