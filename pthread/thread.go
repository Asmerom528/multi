package pthread

/*
#include <pthread.h>

extern void* threadFunc(void*);
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime/cgo"
	"sync"
	"syscall"
	"unsafe"
)

var (
	ErrAlreadyStarted = errors.New("already started")
	ErrNotStarted     = errors.New("not started")
)

// Thread represents a Go function executed in a separate OS thread.
//
// (!) The implementation intentionally ignores the Go runtime safety problems
// of calling Go code from pthread-created threads - use at your own risk.
type Thread struct {
	tid     C.pthread_t
	f       func()
	err     error
	started bool

	mu   sync.Mutex
	once sync.Once
}

// NewThread constructs a Thread ready to be started.
func NewThread(f func()) *Thread {
	return &Thread{f: f}
}

// Start launches the thread. Returns an error if called more than once.
func (t *Thread) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.started {
		return ErrAlreadyStarted
	}

	h := cgo.NewHandle(t) // ensure t is kept alive for the thread lifetime
	//nolint:govet
	ret := C.pthread_create(&t.tid, nil, (*[0]byte)(C.threadFunc), unsafe.Pointer(h))
	if ret != 0 {
		h.Delete()
		return fmt.Errorf("pthread_create: %w", syscall.Errno(ret))
	}

	t.started = true
	return nil
}

// Wait blocks until the thread completes.
// If the thread has already finished, returns immediately.
func (t *Thread) Wait() error {
	t.mu.Lock()
	if !t.started {
		t.mu.Unlock()
		return ErrNotStarted
	}
	t.mu.Unlock()

	t.once.Do(func() {
		ret := C.pthread_join(t.tid, nil)
		if ret != 0 {
			t.err = fmt.Errorf("pthread_join: %w", syscall.Errno(ret))
		} else {
			t.err = nil
		}
	})
	return t.err
}

//export threadFunc
func threadFunc(arg unsafe.Pointer) unsafe.Pointer {
	h := cgo.Handle(arg)
	t := h.Value().(*Thread)
	t.f()
	h.Delete()
	return nil
}
