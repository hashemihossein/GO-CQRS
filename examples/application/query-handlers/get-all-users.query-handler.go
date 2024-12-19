package query_handlers

import (
	"reflect"

	"github.com/hashemihossein/GO-CQRS/examples/application/queries"
	"github.com/hashemihossein/GO-CQRS/examples/domain"
)

type GetAllUsersQueryHandler struct{}

func (handler *GetAllUsersQueryHandler) Handle(query queries.GetAllUsersQuery) ([]domain.User, error) {
	// handling the query, e.g. fetching the users from the database
	return nil, nil
}

func (handler *GetAllUsersQueryHandler) GetQueryType() reflect.Type {
	return reflect.TypeOf(queries.GetAllUsersQuery{})
}
