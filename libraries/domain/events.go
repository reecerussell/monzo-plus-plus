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
	handlers = make(map[string]EventHandler)
)

type EventHandler interface {
	Invoke(ctx context.Context, tx *sql.Tx, e interface{}) errors.Error
}

func RegisterEventHandler(e interface{}, h EventHandler) {
	mu.RLock()
	defer mu.RUnlock()

	noe := reflect.TypeOf(e).Name()

	handlers[noe] = h
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

func (a *Aggregate) RaiseEvent(e interface{}) {
	mu.Lock()
	defer mu.Unlock()

	noe := reflect.TypeOf(e).Name()

	if _, ok := handlers[noe]; !ok {
		panic(fmt.Errorf("no handler registered for event type '%s'", noe))
	}

	a.raisedEvents = append(a.raisedEvents, e)
}

func (a *Aggregate) DispatchEvents(ctx context.Context, tx *sql.Tx) errors.Error {
	mu.Lock()
	defer mu.Unlock()

	for _, e := range a.raisedEvents {
		noe := reflect.TypeOf(e).Name()
		h := handlers[noe]

		err := h.Invoke(ctx, tx, e)
		if err != nil {
			return err
		}
	}

	return nil
}
