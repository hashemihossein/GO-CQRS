package aggregate

import (
	"context"
	"fmt"
	"sync"
)

type AggregateEvent interface {
	Handle() error
}

type AggregateRootInterface interface {
	Apply(event AggregateEvent) error
	ApplyWithoutMiddlewares(event AggregateEvent)
	Commit() error
	CommitWithoutMiddlewares() error
	GetUncommitedEvents() []AggregateEvent
	UseMiddlewareForApply(mw Middleware)
	UseMiddlewareForCommit(mw Middleware)
	LoadFromHistory(events []AggregateEvent) error
}

type AggregateRoot struct {
	events            []AggregateEvent
	applyMiddlewares  []Middleware
	commitMiddlewares []Middleware
	rwMu              sync.RWMutex
}

func (ar *AggregateRoot) Apply(event AggregateEvent) error {
	ar.rwMu.Lock()
	defer ar.rwMu.Unlock()

	chainFunc := func(ctx context.Context, event AggregateEvent) error {
		ar.ApplyWithoutMiddlewares(event)
		return nil
	}

	for i := len(ar.applyMiddlewares) - 1; i >= 0; i-- {
		mw := ar.applyMiddlewares[i]
		next := chainFunc
		chainFunc = func(ctx context.Context, event AggregateEvent) error {
			return mw.Execute(ctx, event, next)
		}
	}

	ctx := context.Background()
	return chainFunc(ctx, event)
}

func (ar *AggregateRoot) ApplyWithoutMiddlewares(event AggregateEvent) {
	ar.events = append(ar.events, event)
}

func (ar *AggregateRoot) Commit() error {
	ar.rwMu.Lock()
	defer ar.rwMu.Unlock()

	for i := 0; i < len(ar.events); i++ {
		err := ar.commitFirstEventWithMiddlewares()
		if err != nil {
			return fmt.Errorf("error while commiting at index %v: %v", i, err.Error())
		}
	}

	ar.events = make([]AggregateEvent, 0)
	return nil
}

func (ar *AggregateRoot) commitFirstEventWithMiddlewares() error {
	event := ar.events[0]

	chainFunc := func(ctx context.Context, event AggregateEvent) error {
		return event.Handle()
	}

	for i := len(ar.commitMiddlewares) - 1; i >= 0; i-- {
		mw := ar.commitMiddlewares[i]
		next := chainFunc
		chainFunc = func(ctx context.Context, event AggregateEvent) error {
			return mw.Execute(ctx, event, next)
		}
	}

	ctx := context.Background()
	return chainFunc(ctx, event)

}

func (ar *AggregateRoot) CommitWithoutMiddlewares() error {
	ar.rwMu.Lock()
	defer ar.rwMu.Unlock()

	for i := 0; i < len(ar.events); i++ {
		event := ar.events[i]
		err := event.Handle()
		if err != nil {
			return fmt.Errorf("error while commiting events at index %v: %v", i, err.Error())
		}
	}

	ar.events = make([]AggregateEvent, 0)
	return nil
}

func (ar *AggregateRoot) UseMiddlewareForApply(mw Middleware) {
	ar.applyMiddlewares = append(ar.applyMiddlewares, mw)
}

func (ar *AggregateRoot) UseMiddlewareForCommit(mw Middleware) {
	ar.commitMiddlewares = append(ar.commitMiddlewares, mw)
}

func (ar *AggregateRoot) LoadFromHistory(events []AggregateEvent) error {
	for i := 0; i < len(events); i++ {
		event := events[i]
		err := event.Handle()
		if err != nil {
			return fmt.Errorf("error while handling events at index %v: %v", i, err.Error())
		}
	}

	return nil
}

func (ar *AggregateRoot) GetUncommitedEvents() []AggregateEvent {
	return ar.events
}

func GetNewAggregateRoot() AggregateRootInterface {
	return &AggregateRoot{
		events:            make([]AggregateEvent, 0),
		applyMiddlewares:  make([]Middleware, 0),
		commitMiddlewares: make([]Middleware, 0),
	}
}
