package query

import "reflect"

type Query interface{}

type QueryResult interface{}

type QueryHandler interface {
	GetQueryType() reflect.Type
}
