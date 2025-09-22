package proc

import "errors"

// Group runs Go functions in separate OS processes.
// Not safe for concurrent use by multiple goroutines.
type Group struct {
	n   int
	res *Chan[string]
	err error
}

// NewGroup creates a new Group.
func NewGroup() *Group {
	return &Group{res: NewChan[string]()}
}

// Go calls the given function in a new OS processes.
func (g *Group) Go(f func() error) {
	p := NewProcess(func() {
		res := f()
		if res != nil {
			g.res.Send(res.Error())
		} else {
			g.res.Send("")
		}
	})

	g.n++
	err := p.Start()
	g.setError(err)
}

// Wait blocks until all processes added with Go have completed,
// then returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	if g.err != nil {
		return g.err
	}
	for range g.n {
		res := g.res.Recv()
		if res != "" {
			g.setError(errors.New(res))
		}
	}
	return g.err
}

// setError sets the group error if it hasn't been set yet.
func (g *Group) setError(err error) {
	if g.err != nil {
		return
	}
	g.err = err
}
