package command

import "reflect"

type Command interface{}

type CommandHandler interface {
	Handle(cmd Command) error
	GetCommandName() reflect.Type
}
