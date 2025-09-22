# Multi: native OS threading and multiprocessing in Go

Multi is a small research project that explores nonconventional ways to handle concurrency in Go by using native OS threads and multiprocessing tools.

It offers three types of "concurrent groups". Each one has an API similar to `sync.WaitGroup`, but they work very differently under the hood:

-   `goro.Group` runs Go functions in goroutines that are locked to OS threads. Each function executes in its own goroutine. Safe to use in production, although unnecessary, because the regular non-locked goroutines work just fine.

-   `pthread.Group` runs Go functions in separate OS threads using POSIX threads. Each function executes in its own thread. This implementation bypasses Go's runtime thread management. Calling Go code from threads not created by the Go runtime can lead to issues with garbage collection, signal handling, and the scheduler. Not meant for production use.

-   `proc.Group` runs Go functions in separate OS processes. Each function executes in its own process forked from the main one. This implementation uses process forking, which is not supported by the Go runtime and can cause undefined behavior, especially in programs with multiple goroutines or complex state. Not meant for production use.

I don't think anyone will find these concurrent groups useful in real-world situations, but it's still interesting to look at a possible (even if flawed) implementations and compare them to Go's default (and only) concurrency model.

## Usage

All groups have a similar API: create a new group, run functions concurrently with `Go`, and wait for completion with `Wait`.

### goro.Group

Runs Go functions in goroutines that are locked to OS threads.

```go
import "github.com/nalgeon/multi/goro"

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
```

You can use channels and other standard concurrency tools inside the functions managed by the group.

### pthread.Group

Runs Go functions in separate OS threads using POSIX threads.

```go
import "github.com/nalgeon/multi/pthread"

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
```

You can use channels and other standard concurrency tools inside the functions managed by the group.

### proc.Group

Runs Go functions in separate OS processes forked from the main one.

```go
import "github.com/nalgeon/multi/proc"

ch := proc.NewChan[int]() // for cross-process communication
defer ch.Close()

g := proc.NewGroup()
g.Go(func() error {
    // do something
    ch.Send(42)
    return nil
})
g.Go(func() error {
    // do something
    ch.Send(42)
    return nil
})
err := g.Wait()
n1, n2 := ch.Recv(), ch.Recv()

fmt.Println(n1, n2, err)
// Output: 42 42 <nil>
```

You can only use `proc.Chan` to exchange data between processes, since regular Go channels and other concurrency tools don't work across process boundaries.

## Benchmarks

Running some CPU-bound workload (with no allocations or I/O) gives these results:

```text
goos: darwin
goarch: arm64
ncpu: 8
gomaxprocs: 8
workers: 4
sync.WaitGroup: n=100   t=60511 µs/exec
goro.Group:     n=100   t=60751 µs/exec
pthread.Group:  n=100   t=60791 µs/exec
proc.Group:     n=100   t=61640 µs/exec
```

One execution here means a group of 4 workers each doing 10 million iterations of generating random numbers and adding them up. See the [benchmark code](internal/benchmark/main.go) for details.

As you can see, the default concurrency model (`sync.WaitGroup` in the results, using standard goroutine scheduling without meddling with threads or processes) works just fine and doesn't add any extra overhead. You probably already knew that, but it's always good to double-check, right?

## Contributing

Contributions are welcome. For anything other than bug fixes, please open an issue first to discuss what you want to change.

Make sure to add or update tests as needed.

## License

Created by [Anton Zhiyanov](https://antonz.org/). Released under the MIT License.
