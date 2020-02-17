package di

import (
	"context"
	"errors"
	"sync"
)

// Service lifetimes.
const (
	// A singleton is created once, then stored in memory for
	// the lifetime of the application. This does not require
	// a context.Context, but can use one to determine variations
	// with in an individual service.
	LifetimeSingleton = 1 << iota

	// A scoped lifetime requires a context.Context, and will stay
	// alive generally for the lifetime of a request, or until the
	// Done() is called on the context. If a scoped service is resolved
	// through the Resolve() method, it will behave as a transient service.
	LifetimeScoped

	// A service with a transient lifetime is built each time Resolve()
	// is called. This service lifetime does not use a context, and will
	// behave the same in each method (Resolve() or ResolveContext()).
	LifetimeTransient
)

var (
	errUnhandledLifetime = errors.New("container: unhandled lifetime")
)

// ServiceLifetime determins the lifetime of the dependency service.
type ServiceLifetime int

// BuilderFunc is a standard function used to build a service.
type BuilderFunc func(ctn *Container) (interface{}, error)

// Service is a struct which defines a service dependeny, which
// can then later be resolved in a Container.
type Service struct {
	Name     string
	Builder  BuilderFunc
	Lifetime ServiceLifetime
}

// Container is a struct that can be used to resolve
// and build service dependencies.
type Container struct {
	mu *sync.RWMutex
	ic *implementationCollection

	// An array of services.
	Services []*Service
}

// New returns a new container instance, with the given service definitions.
func New(services ...*Service) *Container {
	return &Container{
		mu:       &sync.RWMutex{},
		ic:       &implementationCollection{[]*implementation{}},
		Services: services,
	}
}

// Resolve returns a dependency service for the given name. If
// there is no service defined with the given name, nil is returned.
func (ctn *Container) Resolve(name string) interface{} {
	var service *Service

	// Find service.
	for _, s := range ctn.Services {
		if s.Name != name {
			continue
		}

		service = s
	}

	// Return nil, if no service has been found.
	if service == nil {
		return nil
	}

	// Build the service, for the specified lifetime.
	switch service.Lifetime {
	case LifetimeSingleton:
		return ctn.buildSingleton(nil, service)
	case LifetimeScoped:
		// Build scoped service as transient, as
		// no context was given.
		return ctn.buildTransient(service)
	case LifetimeTransient:
		return ctn.buildTransient(service)
	default:
		panic(errUnhandledLifetime)
	}
}

// ResolveContext resolves a dependency service, using the context given.
// The context is only used if the service has a lifetime of scoped.
func (ctn *Container) ResolveContext(ctx context.Context, name string) interface{} {
	var service *Service

	for _, s := range ctn.Services {
		if s.Name != name {
			continue
		}

		service = s
	}

	if service == nil {
		return nil
	}

	switch service.Lifetime {
	case LifetimeSingleton:
		return ctn.buildSingleton(ctx, service)
	case LifetimeScoped:
		return ctn.buildScoped(ctx, service)
	case LifetimeTransient:
		return ctn.buildTransient(service)
	default:
		panic(errUnhandledLifetime)
	}
}

// Clean clears all service definitions and implementations from the container.
// If Clean() has been called on a container, the container can no longer be used.
func (ctn *Container) Clean() {
	ctn.ic = nil
	ctn.Services = nil
	ctn.mu = nil
}

// builds a singleton service. If it has already been built,
// the stored implementation is returned.
func (ctn *Container) buildSingleton(ctx context.Context, s *Service) interface{} {
	i, ok := ctn.ic.Get(ctx, s.Name)
	if ok {
		return i.Build
	}

	return ctn.buildService(ctx, s)
}

// builds a service, using the current context.
func (ctn *Container) buildScoped(ctx context.Context, s *Service) interface{} {
	i, ok := ctn.ic.Get(ctx, s.Name)
	if ok {
		return i.Build
	}

	return ctn.buildService(ctx, s)
}

// builds a service with a transient lifetime.
func (ctn *Container) buildTransient(s *Service) interface{} {
	return ctn.buildService(nil, s)
}

// builds a service using its builder, then stores the implementation.
func (ctn *Container) buildService(ctx context.Context, s *Service) interface{} {
	b, err := s.Builder(ctn)
	if err != nil {
		panic(err)
	}

	if s.Lifetime == LifetimeSingleton ||
		s.Lifetime == LifetimeScoped {

		// Append to implementations.
		ctn.ic.Set(&implementation{
			Name:    s.Name,
			Build:   b,
			Context: ctx,
		})
	}

	return b
}

// implrmentation is a record of a built service,
// which is stored in the container. The identity of an
// implementation is defined by its name and context.
type implementation struct {
	Name    string
	Context context.Context
	Build   interface{}
}

// implementationCollection provides some extension methods
// for a list of implementations.
type implementationCollection struct {
	arr []*implementation
}

// Get attempts to find an implementation in the collection,
// with the same name and context. If one doesn't exist, a nil
// value will be returned, alongside false.
func (ic *implementationCollection) Get(ctx context.Context, name string) (*implementation, bool) {
	for _, i := range ic.arr {
		if i.Name == name {
			if ctx != nil && i.Context != ctx {
				return nil, false
			}

			return i, true
		}
	}

	return nil, false
}

// Set adds an implementation to the collection. Although, if
// an implementation already exists with the same name and context,
// it will be overwritten.
func (ic *implementationCollection) Set(i *implementation) {
	// TODO: remove as the flow dictates that
	// 		 a dependency will never be set if it
	//       already exists.
	// for idx, e := range ic.arr {
	// 	// If an implementation already exists,
	// 	// replace it.
	// 	if e.Name == i.Name {
	// 		if i.Context != nil &&
	// 			e.Context != i.Context {
	// 			continue
	// 		}
	//
	// 		ic.arr[idx] = i
	// 		return
	// 	}
	// }

	// Append a new implementation.
	ic.arr = append(ic.arr, i)
}
