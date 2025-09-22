package multi_test

import (
	"fmt"

	"github.com/nalgeon/multi/goro"
	"github.com/nalgeon/multi/pthread"
)

func Example_goro_Group() {
	ch := make(chan int, 2) // for cross-goroutine communication

	g := goro.NewGroup()
	g.Go(func() error {
		// do something
		ch <- 42
		return nil
	})
	g.Go(func() error {
		// do something
		ch <- 42
		return nil
	})
	err := g.Wait()
	n1, n2 := <-ch, <-ch

	fmt.Println(n1, n2, err)
	// Output: 42 42 <nil>
}

func Example_pthread_Group() {
	ch := make(chan int, 2) // for cross-thread communication

	g := pthread.NewGroup()
	g.Go(func() error {
		// do something
		ch <- 42
		return nil
	})
	g.Go(func() error {
		// do something
		ch <- 42
		return nil
	})
	err := g.Wait()
	n1, n2 := <-ch, <-ch

	fmt.Println(n1, n2, err)
	// Output: 42 42 <nil>
}

// Can't run the proc example in tests due to forked processes.
// See internal/proc/main.go for the example usage.
