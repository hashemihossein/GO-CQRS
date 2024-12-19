package command_handlers

import (
	"reflect"

	"github.com/hashemihossein/GO-CQRS/examples/application/commands"
)

type CreateUserCommandHandler struct{}

func (handler *CreateUserCommandHandler) Handle(cmd commands.CreateUserCommand) error {
	// handling the command, e.g. persisting the user to the database
	return nil
}

func (handler *CreateUserCommandHandler) GetCommandType() reflect.Type {
	return reflect.TypeOf(commands.CreateUserCommand{})
}
