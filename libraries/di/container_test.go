package di

import (
	"context"
	"fmt"
	"testing"
)

// TestService is a basic struct used to test
// service definitions within the container.
type TestService struct {
	Text string
}

// Example test service definitions.
var definitions = []*Service{
	&Service{
		Name:     "one",
		Builder:  getBuilderFunc("a"),
		Lifetime: LifetimeSingleton,
	},
	&Service{
		Name:     "two",
		Builder:  getBuilderFunc("a"),
		Lifetime: LifetimeScoped,
	},
	&Service{
		Name:     "three",
		Builder:  getBuilderFunc("a"),
		Lifetime: LifetimeTransient,
	},
}

// returns a BuilderFunc which builds a TestService with the given text.
func getBuilderFunc(text string) BuilderFunc {
	return func(ctn *Container) (interface{}, error) {
		return &TestService{
			Text: text,
		}, nil
	}
}

func TestNewContainer(t *testing.T) {
	ctn := New(definitions...)

	if len(ctn.Services) != len(definitions) {
		t.Fatalf("expected %d services, but got %d", len(ctn.Services), len(definitions))
	}
}

func TestResolveSingleton(t *testing.T) {
	ctn := New(definitions...)

	// Get initial service.
	s := ctn.Resolve("one").(*TestService)
	if s.Text != "a" {
		t.Fatalf("expected 'a' but got '%s'", s.Text)
	}

	// Change value.
	s.Text = "b"

	// Get the service again, which should have to new value.
	s = ctn.Resolve("one").(*TestService)
	if s.Text != "b" {
		t.Fatalf("expected 'b' but got '%s'", s.Text)
	}
}

func TestResolveSingletonWithContext(t *testing.T) {
	ctn := New(definitions...)
	ctx := context.Background()

	// Get initial service.
	s := ctn.ResolveContext(ctx, "one").(*TestService)
	if s.Text != "a" {
		t.Fatalf("expected 'a' but got '%s'", s.Text)
	}

	// Change value.
	s.Text = "b"

	// Get the service again, which should have to new value.
	s = ctn.ResolveContext(ctx, "one").(*TestService)
	if s.Text != "b" {
		t.Fatalf("expected 'b' but got '%s'", s.Text)
	}
}

func TestResolveScoped(t *testing.T) {
	ctn := New(definitions...)

	// Create a context that is slightly different from
	// the background context. Ignore error.
	ctx := context.WithValue(context.Background(), "test_key", "test_val")

	// Get the initial service.
	s := ctn.ResolveContext(ctx, "two").(*TestService)
	if s.Text != "a" {
		t.Fatalf("expected 'a' but got '%s'", s.Text)
	}

	// Change the value.
	s.Text = "b"

	// Get the service with a different context.
	s = ctn.ResolveContext(context.Background(), "two").(*TestService)
	if s.Text != "a" {
		// The update value will not be applied to this
		// as it has a different context.
		t.Fatalf("expected 'a' but got '%s'", s.Text)
	}

	// Get the initial service again, which should have to new value.
	s = ctn.ResolveContext(ctx, "two").(*TestService)
	if s.Text != "b" {
		t.Fatalf("expected 'b' but got '%s'", s.Text)
	}
}

func TestResolveScopedWithNoContext(t *testing.T) {
	ctn := New(definitions...)

	s1 := ctn.Resolve("two").(*TestService)

	s1.Text = "new text"

	s2 := ctn.Resolve("two").(*TestService)

	if s1.Text == s2.Text {
		t.Fatalf("expected 'a' but got '%s'", s2.Text)
	}
}

func TextResolveTransient(t *testing.T) {
	ctn := New(definitions...)

	// Get the service.
	s := ctn.Resolve("three").(*TestService)
	if s.Text != "a" {
		t.Fatalf("expected 'a' but got '%s'", s.Text)
	}

	// Change value.
	s.Text = "b"

	// Get the service again. This will not have the value
	// as the service has a transient lifetime.
	s = ctn.Resolve("three").(*TestService)
	if s.Text != "a" {
		t.Fatalf("expected 'a' but got '%s'", s.Text)
	}
}

func TestResolveTransientWithContext(t *testing.T) {
	ctn := New(definitions...)
	ctx := context.Background()

	// Get the service.
	s := ctn.ResolveContext(ctx, "three").(*TestService)
	if s.Text != "a" {
		t.Fatalf("expected 'a' but got '%s'", s.Text)
	}
}

func TestBuilderError(t *testing.T) {
	testErr := fmt.Errorf("builder failed")

	s := &Service{
		Name:     "test",
		Lifetime: LifetimeTransient,
		Builder: func(ctn *Container) (interface{}, error) {
			return nil, testErr
		},
	}

	ctn := New(s)

	// Defer a method to recover the panic.
	defer func() {
		if r := recover(); r != nil {
			if r != testErr {
				t.Fatalf("expected '%v' but got '%v'", testErr, r)
			}
		}
	}()

	// Call resolve, which will attempt to build the
	// service, which will result in an error.
	_ = ctn.Resolve("test")
}

func TestResolveNonExistantService(t *testing.T) {
	ctn := New()

	s := ctn.Resolve("non-existant-service")
	if s != nil {
		t.Fatalf("expected nil but got '%v'", s)
	}

	s = ctn.ResolveContext(nil, "non-existant-service")
	if s != nil {
		t.Fatalf("expected nil but got '%v'", s)
	}
}

func TestUnhandledLifetime(t *testing.T) {
	s := &Service{
		Name:     "e",
		Lifetime: 999, // an unhandled lifetime
	}

	// Create the container with the serive def.
	ctn := New(s)

	defer func() {
		if r := recover(); r != nil {
			if r != errUnhandledLifetime {
				t.Fatalf("expected '%v' but got '%v'", errUnhandledLifetime, r)
			}
		}
	}()

	_ = ctn.Resolve("e")
}

func TestUnhandledLifetimeWithContext(t *testing.T) {
	s := &Service{
		Name:     "e",
		Lifetime: 999, // an unhandled lifetime
	}

	// Create the container with the serive def.
	ctn := New(s)

	defer func() {
		if r := recover(); r != nil {
			if r != errUnhandledLifetime {
				t.Fatalf("expected '%v' but got '%v'", errUnhandledLifetime, r)
			}
		}
	}()

	_ = ctn.ResolveContext(context.Background(), "e")
}
