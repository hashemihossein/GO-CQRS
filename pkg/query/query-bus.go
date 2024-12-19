package query

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type QueryBus struct {
	handlers    map[reflect.Type]func(Query) (QueryResult, error)
	middlewares []Middleware
	rwMu        sync.RWMutex
}

type QueryBusInterface interface {
	RegisterQueryHandler(handler QueryHandler) error
	Dispatch(query Query) (QueryResult, error)
	DispatchWithContext(ctx context.Context, query Query) (QueryResult, error)
	DispatchWithoutMiddlewares(query Query) (QueryResult, error)
	UseMiddleware(mw Middleware)
}

var queryBusInstance *QueryBus
var once sync.Once

func GetQueryBus() QueryBusInterface {
	once.Do(func() {
		queryBusInstance = &QueryBus{
			handlers:    make(map[reflect.Type]func(Query) (QueryResult, error)),
			middlewares: make([]Middleware, 0),
		}
	})

	return queryBusInstance
}

func (qb *QueryBus) RegisterQueryHandler(handler QueryHandler) error {
	qb.rwMu.Lock()
	defer qb.rwMu.Unlock()

	queryType := handler.GetQueryType()
	if _, exists := qb.handlers[queryType]; exists {
		return errors.New("this query handler has been registered before")
	}

	handlerValue := reflect.ValueOf(handler)
	handlerType := handlerValue.Type()

	method, exists := handlerType.MethodByName("Handle")
	if !exists {
		return errors.New("handler does not have a Handle method")
	}

	if method.Type.NumIn() != 2 || method.Type.NumOut() != 2 {
		return errors.New("Handle method has incorrect signature, expected Handle(SpecificQuery) (QueryResult, error)")
	}

	specificQueryType := method.Type.In(1)
	if specificQueryType != queryType {
		return errors.New("Handle method parameter type does not match GetQueryType return type")
	}

	expectedErrorType := reflect.TypeOf((*error)(nil)).Elem()
	if method.Type.Out(1) != expectedErrorType {
		return errors.New("Handle method must return error as the second return value")
	}

	wrapper := func(query Query) (QueryResult, error) {
		if reflect.TypeOf(query) != queryType {
			return nil, errors.New("invalid query type")
		}

		results := handlerValue.MethodByName("Handle").Call([]reflect.Value{reflect.ValueOf(query)})

		if len(results) != 2 {
			return nil, errors.New("Handle method should return exactly two values")
		}

		var result QueryResult
		if results[0].CanInterface() {
			result = results[0].Interface().(QueryResult)
		}

		var err error
		if results[1].CanInterface() && results[1].Interface() != nil {
			err = results[1].Interface().(error)
		}

		return result, err
	}

	qb.handlers[queryType] = wrapper
	return nil
}

func (qb *QueryBus) Dispatch(query Query) (QueryResult, error) {
	qb.rwMu.RLock()
	_, exists := qb.handlers[reflect.TypeOf(query)]
	qb.rwMu.RUnlock()
	if !exists {
		return nil, errors.New("there is no registered handler for this query")
	}

	ctx := context.Background()
	return qb.DispatchWithContext(ctx, query)
}

func (qb *QueryBus) DispatchWithContext(ctx context.Context, query Query) (QueryResult, error) {
	qb.rwMu.RLock()
	handlerFunc, exists := qb.handlers[reflect.TypeOf(query)]
	if !exists {
		qb.rwMu.RUnlock()
		return nil, errors.New("there is no registered handler for this query")
	}

	middlewares := make([]Middleware, len(qb.middlewares))
	copy(middlewares, qb.middlewares)
	qb.rwMu.RUnlock()

	chainFunc := func(ctx context.Context, query Query) (QueryResult, error) {
		return handlerFunc(query)
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		next := chainFunc
		chainFunc = func(ctx context.Context, query Query) (QueryResult, error) {
			return mw.Execute(ctx, query, next)
		}
	}

	return chainFunc(ctx, query)
}

func (qb *QueryBus) DispatchWithoutMiddlewares(query Query) (QueryResult, error) {
	qb.rwMu.RLock()
	handlerFunc, exists := qb.handlers[reflect.TypeOf(query)]
	qb.rwMu.RUnlock()
	if !exists {
		return nil, errors.New("there is no registered handler for this query")
	}

	return handlerFunc(query)
}

func (qb *QueryBus) UseMiddleware(mw Middleware) {
	qb.rwMu.Lock()
	defer qb.rwMu.Unlock()
	qb.middlewares = append(qb.middlewares, mw)
}
