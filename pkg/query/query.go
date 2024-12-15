package query

import "reflect"

type Query interface{}

type QueryResult interface{}

type QueryHandler interface {
	Handle(query Query) (QueryResult, error)
	GetQueryName() reflect.Type
}
