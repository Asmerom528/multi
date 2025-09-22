// Package pthread provides a concurrency group that runs functions
// in separate OS threads using POSIX threads.
//
// Usage example:
//
//	g := pthread.NewGroup()
//	for range 8 {
//	    g.Go(func() error {
//	        // do some work here
//	        return nil
//	    })
//	}
//	if err := g.Wait(); err != nil {
//	    // handle error
//	}
package pthread
