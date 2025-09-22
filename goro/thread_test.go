package goro

import (
	"testing"
	"time"
)

func TestThread(t *testing.T) {
	t.Run("new thread", func(t *testing.T) {
		th := NewThread(func() {})
		if th == nil {
			t.Fatal("returned nil")
		}
		if th.started {
			t.Error("should not start the thread")
		}
	})
	t.Run("start", func(t *testing.T) {
		th := NewThread(func() {})
		err := th.Start()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if !th.started {
			t.Error("started flag not set")
		}
	})
	t.Run("multiple start", func(t *testing.T) {
		th := NewThread(func() {})
		_ = th.Start()
		err := th.Start()
		if err != ErrAlreadyStarted {
			t.Errorf("got %v, want %v", err, ErrAlreadyStarted)
		}
	})
	t.Run("wait before start", func(t *testing.T) {
		th := NewThread(func() {})
		err := th.Wait()
		if err != ErrNotStarted {
			t.Errorf("got %v, want %v", err, ErrNotStarted)
		}
	})
	t.Run("start and wait", func(t *testing.T) {
		var done bool
		th := NewThread(func() { done = true })
		_ = th.Start()
		err := th.Wait()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if !done {
			t.Error("function was not executed")
		}
	})
	t.Run("multiple wait", func(t *testing.T) {
		th := NewThread(func() {})
		_ = th.Start()
		_ = th.Wait()
		err := th.Wait()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
	})
	t.Run("concurrent wait", func(t *testing.T) {
		th := NewThread(func() {
			time.Sleep(10 * time.Millisecond)
		})
		_ = th.Start()
		done := make(chan bool, 2)
		for range 2 {
			go func() {
				_ = th.Wait()
				done <- true
			}()
		}
		for range 2 {
			select {
			case <-done:
			case <-time.After(50 * time.Millisecond):
				t.Error("Wait did not complete in time")
			}
		}
	})
}
