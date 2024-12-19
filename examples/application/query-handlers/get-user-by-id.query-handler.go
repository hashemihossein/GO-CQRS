package query_handlers

import (
	"reflect"

	"github.com/hashemihossein/GO-CQRS/examples/application/queries"
	"github.com/hashemihossein/GO-CQRS/examples/domain"
	domainEvents "github.com/hashemihossein/GO-CQRS/examples/domain/events"
)

type GetUserByIdQueryHandler struct{}

func (handler *GetUserByIdQueryHandler) Handle(query queries.GetUserByIdQuery) (domain.User, error) {
	// handling the query, e.g. fetching the user from the database
	user := domain.NewUser("username", "password", "2000-01-01")
	user.Apply(&domainEvents.UserCreatedEvent{User: user})
	user.Commit()
	return *user, nil
}

func (handler *GetUserByIdQueryHandler) GetQueryType() reflect.Type {
	return reflect.TypeOf(queries.GetUserByIdQuery{})
}
