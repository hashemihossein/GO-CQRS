package command

import "reflect"

type Command interface{}

type CommandHandler interface {
	GetCommandType() reflect.Type
}
