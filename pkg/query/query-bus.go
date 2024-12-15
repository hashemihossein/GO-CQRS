package query

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type QueryBusInterface interface {
	RegisterQueryHandler(handler QueryHandler) error
	Dispatch(query Query) (QueryResult, error)
	DispatchWithContext(ctx context.Context, query Query) (QueryResult, error)
	DispatchWithoutMiddlewares(query Query) (QueryResult, error)
	UseMiddleware(mw QueryMiddleware)
}

type QueryBus struct {
	queryHandlers map[reflect.Type]QueryHandler
	middlewares   []QueryMiddleware
	rwMu          sync.RWMutex
}

var queryBusInstance *QueryBus

func (qb *QueryBus) RegisterQueryHandler(handler QueryHandler) error {
	queryName := handler.GetQueryName()
	handler, exists := qb.queryHandlers[queryName]
	if exists {
		errors.New("this query handler have been registered before")
	}

	qb.queryHandlers[queryName] = handler

	return nil
}

func (qb *QueryBus) Dispatch(query Query) (QueryResult, error) {
	qb.rwMu.RLock()

	queryType := reflect.TypeOf(query)
	_, exists := qb.queryHandlers[queryType]
	if !exists {
		return nil, errors.New("there is not registered query handler for this query")
	}

	qb.rwMu.RUnlock()

	ctx := context.Background()
	return qb.DispatchWithContext(ctx, query)
}

func (qb *QueryBus) DispatchWithContext(ctx context.Context, query Query) (QueryResult, error) {
	qb.rwMu.RLock()
	defer qb.rwMu.RUnlock()

	queryType := reflect.TypeOf(query)
	handler, exists := qb.queryHandlers[queryType]
	if !exists {
		return nil, errors.New("there is not registered query handler for this query")
	}

	chainFunc := func(ctx context.Context, query Query) (QueryResult, error) {
		return handler.Handle(query)
	}

	for i := len(qb.middlewares); i >= 0; i-- {
		mw := qb.middlewares[i]
		next := chainFunc
		chainFunc = func(ctx context.Context, query Query) (QueryResult, error) {
			return mw.Execute(ctx, query, next)
		}
	}

	return chainFunc(ctx, query)
}

func (qb *QueryBus) DispatchWithoutMiddlewares(query Query) (QueryResult, error) {
	qb.rwMu.RLock()
	defer qb.rwMu.RUnlock()

	queryType := reflect.TypeOf(query)
	handler, exists := qb.queryHandlers[queryType]
	if !exists {
		return nil, errors.New("there is not registered query handler for this query")
	}

	return handler.Handle(query)
}

func (qb *QueryBus) UseMiddleware(mw QueryMiddleware) {
	qb.middlewares = append(qb.middlewares, mw)
}

var once sync.Once

func getQueryBus() QueryBusInterface {
	once.Do(func() {
		queryBusInstance = &QueryBus{
			queryHandlers: make(map[reflect.Type]QueryHandler),
			middlewares:   make([]QueryMiddleware, 0),
		}
	})

	return queryBusInstance
}
