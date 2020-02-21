package errors

import "sync"

type GroupFunc func() Error

// Group is used to handle multiple errors concurrently.
type Group struct {
	wg sync.WaitGroup

	errOnce sync.Once
	err     Error
}

// Go appends a GroupFunc to the groups function stack.
func (g *Group) Go(f GroupFunc) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
	}()
}

// Wait executes each of the group's functions, in parallel.
// The first err, if any, will be the error in which is returned.
func (g *Group) Wait() Error {
	g.wg.Wait()

	return g.err
}
