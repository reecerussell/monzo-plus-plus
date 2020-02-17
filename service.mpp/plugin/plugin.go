package plugin

import (
	"context"
	"net/http"
	"sync"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/monzo"
)

// TransactionFunc is a function type that is called on an event.
type TransactionFunc func(ctx context.Context, t *monzo.Transaction) error

// Plugin is used to define and consume plugins.
type Plugin interface {
	// Name is used for informational purposes, to document a plugin.
	// Both Name() and Description() are used to explain what a plugin
	// is and what it is doing.
	//
	// These don't have to return anything, but should as it makes the plugin
	// more accessible.
	Name() string
	Description() string

	// Build gives a plugin the change to use the service repository
	// to use services such as the user usecase.
	//
	// A plugin must implement this method, but does not have to make
	// use of it.
	Build(ctn *di.Container)

	// Handler is used to provide an end user or developer with a set of HTTP
	// endpoints as a RESTful API. This is entirely optional, but could be useful
	// in some cases.
	//
	// This field can return a nil value.
	Handler() http.Handler

	// TransactionCreate will be called each time a transaction is created. This method
	// provides the plugin an entrypoint that can be used to execute an operation
	// on this event.
	TransactionCreated(ctx context.Context, t *monzo.Transaction) error
}

var (
	mu      = &sync.RWMutex{}
	plugins = []Plugin{}
)

// Register registers a new plugin that will be used on an event.
func Register(p Plugin) {
	mu.Lock()
	defer mu.Unlock()

	plugins = append(plugins, p)
}

// Build call Build() on each plugin, passing the di.Container.
func Build(ctn *di.Container) {
	mu.RLock()
	defer mu.RUnlock()

	for _, p := range plugins {
		p.Build(ctn)
	}
}

// TransactionCreatedHandler returns all handler methods for the transaction created event.
func TransactionCreatedHandler() []TransactionFunc {
	mu.RLock()
	defer mu.RUnlock()

	h := make([]TransactionFunc, len(plugins))

	for i, f := range plugins {
		h[i] = f.TransactionCreated
	}

	return h
}
