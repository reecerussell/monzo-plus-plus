package errors

import (
	"testing"
	"time"
)

func TestGroupError(t *testing.T) {
	g := new(Group)

	e1 := NotFound("Err 1")
	e2 := NotFound("Err 2")

	g.Go(func() Error {
		return nil
	})
	g.Go(func() Error {
		time.Sleep(1 * time.Second)
		return e1
	})
	g.Go(func() Error {
		time.Sleep(2 * time.Second)
		return e2
	})

	if err := g.Wait(); err != e1 {
		t.Fatalf("Expected: %s, but got %s", e1.Text(), err.Text())
	}
}

func TestGroupNoError(t *testing.T) {
	g := new(Group)

	g.Go(func() Error {
		return nil
	})
	g.Go(func() Error {
		return nil
	})

	if err := g.Wait(); err != nil {
		t.Fatalf("Expected: nil, but got %s", err.Text())
	}
}
