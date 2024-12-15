package query

import "context"

type QueryMiddleware interface {
	Execute(ctx context.Context, query Query, next NextFunc) (QueryResult, error)
}

type NextFunc func(ctx context.Context, query Query) (QueryResult, error)
