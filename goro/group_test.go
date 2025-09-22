package goro

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	t.Run("new group", func(t *testing.T) {
		g := NewGroup()
		if g == nil {
			t.Fatal("NewGroup returned nil")
		}
	})
	t.Run("go", func(t *testing.T) {
		g := NewGroup()
		called := false
		g.Go(func() error {
			called = true
			return nil
		})
		err := g.Wait()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if !called {
			t.Error("function was not called")
		}
	})
	t.Run("multiple go", func(t *testing.T) {
		g := NewGroup()
		var count atomic.Int32
		for range 5 {
			g.Go(func() error {
				time.Sleep(10 * time.Millisecond) // Simulate work
				count.Add(1)
				return nil
			})
		}
		err := g.Wait()
		if err != nil {
			t.Errorf("got error %v, want nil", err)
		}
		if count.Load() != 5 {
			t.Errorf("got count %d, want 5", count.Load())
		}
	})
	t.Run("error", func(t *testing.T) {
		g := NewGroup()
		testErr := errors.New("test error")
		g.Go(func() error {
			return testErr
		})
		g.Go(func() error {
			return nil
		})
		err := g.Wait()
		if err != testErr {
			t.Errorf("got %v, want %v", err, testErr)
		}
	})
	t.Run("multiple errors", func(t *testing.T) {
		g := NewGroup()
		firstErr := errors.New("first")
		secondErr := errors.New("second")
		g.Go(func() error {
			time.Sleep(20 * time.Millisecond)
			return firstErr
		})
		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			return secondErr
		})
		err := g.Wait()
		if err != secondErr {
			t.Errorf("got %v, want %v (first error to occur)", err, secondErr)
		}
	})
}
