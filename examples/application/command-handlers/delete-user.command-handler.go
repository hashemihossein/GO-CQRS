package command_handlers

import (
	"reflect"

	"github.com/hashemihossein/GO-CQRS/examples/application/commands"
)

type DeleteUserCommandHandler struct{}

func (handler *DeleteUserCommandHandler) Handle(cmd commands.DeleteUserCommand) error {
	// handling the command, e.g. deleting the user from the database
	return nil
}

func (handler *DeleteUserCommandHandler) GetCommandType() reflect.Type {
	return reflect.TypeOf(commands.DeleteUserCommand{})
}
