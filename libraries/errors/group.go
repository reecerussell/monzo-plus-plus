package errors

type GroupFunc func() Error

// Group is used to handle multiple errors concurrently.
type Group struct {
	fs []GroupFunc
}

// Go appends a GroupFunc to the groups function stack.
func (g *Group) Go(f GroupFunc) {
	g.fs = append(g.fs, f)
}

// Wait executes each of the group's functions, in parallel.
// The first err, if any, will be the error in which is returned.
func (g *Group) Wait() Error {
	errs := make(chan Error)

	for _, f := range g.fs {
		go func(gf GroupFunc) {
			errs <- gf()
		}(f)
	}

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
