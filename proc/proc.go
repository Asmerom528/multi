// Package proc provides a concurrency group that
// runs functions in separate OS processes.
//
// Usage example:
//
//	g := proc.NewGroup()
//	for range 8 {
//	    g.Go(func() error {
//	        // do some work here
//	        return nil
//	    })
//	}
//	if err := g.Wait(); err != nil {
//	    // handle error
//	}
package proc
