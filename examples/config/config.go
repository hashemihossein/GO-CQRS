package config

import (
	"fmt"
	"log"

	commandHandlers "github.com/hashemihossein/GO-CQRS/examples/application/command-handlers"
	queryHandlers "github.com/hashemihossein/GO-CQRS/examples/application/query-handlers"
	"github.com/hashemihossein/GO-CQRS/pkg/command"
	"github.com/hashemihossein/GO-CQRS/pkg/query"
)

func checkErrors(errs []error) error {
	for _, err := range errs {
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func registerCommandHandlers() error {
	fmt.Println("Registering command handlers")
	commandBus := command.GetCommandBus()

	errs := []error{}

	errs = append(errs, commandBus.RegisterCommandHandler(&commandHandlers.CreateUserCommandHandler{}))
	errs = append(errs, commandBus.RegisterCommandHandler(&commandHandlers.DeleteUserCommandHandler{}))

	return checkErrors(errs)
}

func registerQueryHandlers() error {
	fmt.Println("Registering query handlers")
	queryBus := query.GetQueryBus()

	errs := []error{}

	errs = append(errs, queryBus.RegisterQueryHandler(&queryHandlers.GetAllUsersQueryHandler{}))
	errs = append(errs, queryBus.RegisterQueryHandler(&queryHandlers.GetUserByIdQueryHandler{}))

	return checkErrors(errs)
}

func Register() {
	registerCommandHandlers()
	registerQueryHandlers()
}
