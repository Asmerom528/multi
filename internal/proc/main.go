// Example usage and test cases for the proc package.
// Can't be in tests due to forked processes.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nalgeon/multi/proc"
)

func example() {
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
}

func cases() {
	var exitCode int
	exitCode |= runCase("new process", func() error {
		p := proc.NewProcess(func() {})
		if p == nil {
			return fmt.Errorf("returned nil")
		}
		return nil
	})
	exitCode |= runCase("start", func() error {
		p := proc.NewProcess(func() {})
		err := p.Start()
		if err != nil {
			return fmt.Errorf("got error %v, want nil", err)
		}
		return nil
	})
	exitCode |= runCase("multiple start", func() error {
		p := proc.NewProcess(func() {})
		_ = p.Start()
		err := p.Start()
		if err != proc.ErrAlreadyStarted {
			return fmt.Errorf("got %v, want %v", err, proc.ErrAlreadyStarted)
		}
		return nil
	})
	exitCode |= runCase("wait before start", func() error {
		p := proc.NewProcess(func() {})
		err := p.Wait()
		if err != proc.ErrNotStarted {
			return fmt.Errorf("got %v, want %v", err, proc.ErrNotStarted)
		}
		return nil
	})
	exitCode |= runCase("start and wait", func() error {
		ch := proc.NewChan[bool]()
		defer ch.Close()

		p := proc.NewProcess(func() {
			ch.Send(true)
		})

		err := p.Start()
		if err != nil {
			return fmt.Errorf("Start: %v", err)
		}

		err = p.Wait()
		if err != nil {
			return fmt.Errorf("Wait: got error %v, want nil", err)
		}

		time.Sleep(50 * time.Millisecond)
		done := ch.Recv()
		if !done {
			return fmt.Errorf("function was not executed")
		}
		return nil
	})
	exitCode |= runCase("multiple wait", func() error {
		p := proc.NewProcess(func() {})
		_ = p.Start()
		_ = p.Wait()
		err := p.Wait()
		if err != nil {
			return fmt.Errorf("got error %v, want nil", err)
		}
		return nil
	})
	exitCode |= runCase("concurrent wait", func() error {
		p := proc.NewProcess(func() {
			time.Sleep(10 * time.Millisecond)
		})
		_ = p.Start()
		done := make(chan bool, 2)
		for range 2 {
			go func() {
				_ = p.Wait()
				done <- true
			}()
		}
		for range 2 {
			select {
			case <-done:
			case <-time.After(50 * time.Millisecond):
				return fmt.Errorf("Wait did not complete in time")
			}
		}
		return nil
	})
	if exitCode != 0 {
		fmt.Println("FAIL")
		os.Exit(exitCode)
	}
	fmt.Println("PASS")
	os.Exit(exitCode)
}

func runCase(name string, f func() error) int {
	err := f()
	if err != nil {
		fmt.Printf("✗ %v... ERROR: %v\n", name, err)
		return 1
	}
	fmt.Printf("✓ %v... OK\n", name)
	return 0
}

func main() {
	example()
	cases()
}
