package proc

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

var (
	ErrAlreadyStarted = errors.New("already started")
	ErrNotStarted     = errors.New("not started")
)

// Process represents a Go function executed in a forked OS process.
// Not safe for concurrent use by multiple goroutines.
//
// (!) The implementation intentionally ignores the Go runtime safety problems
// of calling Go code from forked processes - use at your own risk.
type Process struct {
	f       func()
	pid     int
	err     error
	started bool
	done    bool
}

// NewProcess constructs a Process ready to be started.
func NewProcess(f func()) *Process {
	return &Process{f: f}
}

// Start launches the process. Returns an error if called more than once.
func (p *Process) Start() error {
	if p.started {
		return ErrAlreadyStarted
	}
	p.started = true

	pid, r2, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return fmt.Errorf("syscall.Fork: %w", err)
	}

	if runtime.GOOS == "darwin" && r2 == 1 {
		// On Darwin:
		//	pid = child pid in both parent and child.
		//	r2 = 0 in parent, 1 in child.
		// Convert to normal Unix pid = 0 in child.
		// https://github.com/golang/go/blob/go1.11/src/syscall/exec_bsd.go#L76
		pid = 0
	}

	if pid != 0 {
		// In parent.
		p.pid = int(pid)
		return nil
	}

	// Fork succeeded, now in child.
	p.f()
	os.Exit(0)
	return nil
}

// Wait blocks until the process completes.
// If the process has already finished, returns immediately.
func (p *Process) Wait() error {
	if !p.started {
		return ErrNotStarted
	}
	if p.done {
		return p.err
	}

	var ws syscall.WaitStatus
	_, err := syscall.Wait4(p.pid, &ws, 0, nil)
	if err != nil {
		p.err = fmt.Errorf("syscall.Wait4: %w", err)
		return p.err
	}
	if !ws.Exited() || ws.ExitStatus() != 0 {
		p.err = fmt.Errorf("exit status %d", ws.ExitStatus())
		return p.err
	}

	p.done = true
	return nil
}
