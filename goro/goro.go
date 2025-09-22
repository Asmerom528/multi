// Package goro provides a concurrency group that runs functions
// in separate goroutines, each locked to its own OS thread.
//
// Usage example:
//
//	g := goro.NewGroup()
//	for range 8 {
//	    g.Go(func() error {
//	        // do some work here
//	        return nil
//	    })
//	}
//	if err := g.Wait(); err != nil {
//	    // handle error
//	}
package goro
