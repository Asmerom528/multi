// Tracing different concurrency group implementations.
// Compares sync.WaitGroup, goro.Group, pthread.Group, and proc.Group.
package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path"
	"runtime"
	"runtime/trace"
	"sync"

	"github.com/nalgeon/multi/goro"
	"github.com/nalgeon/multi/proc"
	"github.com/nalgeon/multi/pthread"
)

const nWorkers = 4
const nIter = 10_000_000

var sink any

// concGroup is a common interface for different concurrency group implementations.
type concGroup interface {
	Go(func() error)
	Wait() error
}

// syncGroup is a wrapper around sync.WaitGroup to implement concGroup interface.
type syncGroup struct {
	wg sync.WaitGroup
}

func (s *syncGroup) Go(f func() error) {
	s.wg.Go(func() { _ = f() })
}

func (s *syncGroup) Wait() error {
	s.wg.Wait()
	return nil
}

// work simulates some work by generating random numbers and summing them up.
func work() error {
	var total int
	for range nIter {
		n := rand.N(100)
		total += n
	}
	sink = total
	return nil
}

// collect collects trace for a given concurrency group implementation.
// Creates a new group with nWorkers, starts work in each worker, and waits
// for all workers complete.
func collect(name string, fnNewGroup func() concGroup, work func() error) {
	fname := path.Join("build", name+".trace")
	f, _ := os.Create(fname)
	_ = trace.Start(f)
	defer trace.Stop()

	g := fnNewGroup()
	for range nWorkers {
		g.Go(work)
	}
	err := g.Wait()
	if err != nil {
		log.Fatalf("%s: got error %v, want nil", name, err)
	}

	fmt.Printf("✓ %v: %v\n", name, fname)
}

func main() {
	fmt.Println("Tracing concurrency group implementations...")
	fmt.Println("goos:", runtime.GOOS)
	fmt.Println("goarch:", runtime.GOARCH)
	fmt.Println("ncpu:", runtime.NumCPU())
	fmt.Println("gomaxprocs:", runtime.GOMAXPROCS(0))
	fmt.Println("workers:", nWorkers)
	{
		fnNew := func() concGroup { return &syncGroup{} }
		collect("sync.WaitGroup", fnNew, work)
	}
	{
		fnNew := func() concGroup { return goro.NewGroup() }
		collect("goro.Group", fnNew, work)
	}
	{
		fnNew := func() concGroup { return pthread.NewGroup() }
		collect("pthread.Group", fnNew, work)
	}
	{
		fnNew := func() concGroup { return proc.NewGroup() }
		collect("proc.Group", fnNew, work)
	}
	fmt.Println("Done.")
	_ = sink
}
