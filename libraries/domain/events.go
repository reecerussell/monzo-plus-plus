package domain

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

var (
	mu       = sync.RWMutex{}
	handlers = make(map[reflect.Type]EventHandler)
)

type Event interface{}

type EventHandler interface {
	Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error
}

func RegisterEventHandler(e Event, h EventHandler) {
	mu.RLock()
	defer mu.RUnlock()

	handlers[reflect.TypeOf(e)] = h
}

type Aggregate struct {
	raisedEvents []interface{}
}

func (a *Aggregate) GetRaisedEvents() []interface{} {
	if a.raisedEvents == nil {
		return []interface{}{}
	}

	return a.raisedEvents
}

func (a *Aggregate) RaiseEvent(e Event) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := handlers[reflect.TypeOf(e)]; !ok {
		panic(fmt.Errorf("no handler registered for event type '%s'", reflect.TypeOf(e)))
	}

	a.raisedEvents = append(a.raisedEvents, e)
}

func (a *Aggregate) DispatchEvents(ctx context.Context, tx *sql.Tx) errors.Error {
	mu.Lock()
	defer mu.Unlock()

	for _, e := range a.raisedEvents {
		h := handlers[reflect.TypeOf(e)]

		err := h.Invoke(ctx, tx, e)
		if err != nil {
			return err
		}
	}

	return nil
}
