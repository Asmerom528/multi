package goro

import (
	"errors"
	"runtime"
	"sync"
)

var (
	ErrAlreadyStarted = errors.New("already started")
	ErrNotStarted     = errors.New("not started")
)

// Thread represents a goroutine locked to an OS thread.
type Thread struct {
	f       func()
	done    chan struct{}
	started bool
	mu      sync.Mutex
}

// NewThread constructs a Thread ready to be started.
func NewThread(f func()) *Thread {
	return &Thread{
		f:    f,
		done: make(chan struct{}),
	}
}

// Start launches the goroutine. Returns an error if called more than once.
func (t *Thread) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.started {
		return ErrAlreadyStarted
	}

	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		t.f()
		close(t.done)
	}()

	t.started = true
	return nil
}

// Wait blocks until the goroutine completes.
// If the goroutine has already finished, returns immediately.
func (t *Thread) Wait() error {
	t.mu.Lock()
	if !t.started {
		t.mu.Unlock()
		return ErrNotStarted
	}
	t.mu.Unlock()

	<-t.done
	return nil
}
