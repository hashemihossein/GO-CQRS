package main

import (
	"github.com/hashemihossein/GO-CQRS/examples/application/commands"
	"github.com/hashemihossein/GO-CQRS/examples/application/queries"
	"github.com/hashemihossein/GO-CQRS/examples/config"
	"github.com/hashemihossein/GO-CQRS/pkg/command"
	"github.com/hashemihossein/GO-CQRS/pkg/query"
)

func main() {
	config.Register()
	commandBus := command.GetCommandBus()
	queryBus := query.GetQueryBus()

	commandBus.Dispatch(commands.CreateUserCommand{Username: "test_username", Password: "test_password", DateOfBirth: "test_date_of_birth"})
	queryBus.Dispatch(queries.GetUserByIdQuery{ID: "test_id"})
}
