// Benchmarking different concurrency group implementations.
// Compares sync.WaitGroup, goro.Group, pthread.Group, and proc.Group.
// Measures the average time taken to execute a fixed number of tasks
// across a fixed number of workers.
package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"

	"github.com/nalgeon/multi/goro"
	"github.com/nalgeon/multi/proc"
	"github.com/nalgeon/multi/pthread"
)

const nTimes = 100
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

// benchmark runs the benchmark for a given concurrency group implementation.
// Creates a new group with nWorkers, starts work in each worker, and waits
// for all workers complete. Measures the average time per group execution
// over nTimes iterations.
func benchmark(name string, fnNewGroup func() concGroup, work func() error, nTimes int) {
	start := time.Now()
	for range nTimes {
		g := fnNewGroup()
		for range nWorkers {
			g.Go(work)
		}
		err := g.Wait()
		if err != nil {
			log.Fatalf("%s: got error %v, want nil", name, err)
		}
	}

	totalNs := time.Since(start).Nanoseconds()
	avgDur := time.Duration(totalNs / int64(nTimes))
	fmt.Printf("%v:\tn=%v\tt=%v µs/exec\n", name, nTimes, avgDur.Microseconds())
}

func main() {
	fmt.Println("Benchmarking concurrency group implementations...")
	fmt.Println("goos:", runtime.GOOS)
	fmt.Println("goarch:", runtime.GOARCH)
	fmt.Println("ncpu:", runtime.NumCPU())
	fmt.Println("gomaxprocs:", runtime.GOMAXPROCS(0))
	fmt.Println("workers:", nWorkers)
	{
		fnNew := func() concGroup { return &syncGroup{} }
		benchmark("sync.WaitGroup", fnNew, work, nTimes)
	}
	{
		fnNew := func() concGroup { return goro.NewGroup() }
		benchmark("goro.Group", fnNew, work, nTimes)
	}
	{
		fnNew := func() concGroup { return pthread.NewGroup() }
		benchmark("pthread.Group", fnNew, work, nTimes)
	}
	{
		fnNew := func() concGroup { return proc.NewGroup() }
		benchmark("proc.Group", fnNew, work, nTimes)
	}
	fmt.Println("Done.")
	_ = sink
}
