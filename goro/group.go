package goro

import "sync"

// Group runs Go functions in separate goroutines.
// Each goroutine is locked to its own OS thread.
type Group struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
}

// NewGroup creates a new Group.
func NewGroup() *Group {
	return &Group{}
}

// Go calls the given function in a new thread-locked goroutine.
func (g *Group) Go(f func() error) {
	t := NewThread(func() {
		defer g.wg.Done()
		err := f()
		if err != nil {
			g.setError(err)
		}
	})
	g.wg.Add(1)
	err := t.Start()
	if err != nil {
		g.setError(err)
		g.wg.Done()
	}
}

// Wait blocks until all goroutines added with Go have completed,
// then returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}

// setError sets the group error if it hasn't been set yet.
func (g *Group) setError(err error) {
	g.errOnce.Do(func() {
		g.err = err
	})
}
