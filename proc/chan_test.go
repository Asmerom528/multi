package proc

import (
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		ch := NewChan[int]()
		defer ch.Close()

		go func() {
			ch.Send(42)
		}()

		got := ch.Recv()
		want := 42
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("string", func(t *testing.T) {
		ch := NewChan[string]()
		defer ch.Close()

		go func() {
			ch.Send("hello")
		}()

		got := ch.Recv()
		want := "hello"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("multiple send-recv", func(t *testing.T) {
		ch := NewChan[int]()
		defer ch.Close()

		go func() {
			for i := range 3 {
				ch.Send(i)
			}
		}()

		for i := 0; i < 3; i++ {
			got := ch.Recv()
			want := i
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})
	t.Run("recv then send", func(t *testing.T) {
		ch := NewChan[int]()
		defer ch.Close()

		go func() {
			time.Sleep(10 * time.Millisecond)
			ch.Send(100)
		}()

		got := ch.Recv()
		want := 100
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("recv after close", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Error("should not panic on Recv after Close")
			}
		}()

		ch := NewChan[int]()
		ch.Close()
		val := ch.Recv()
		if val != 0 {
			t.Errorf("got %v, want 0", val)
		}
	})
	t.Run("close", func(t *testing.T) {
		ch := NewChan[int]()
		ch.Close()
		if !ch.Closed() {
			t.Error("expected channel to be closed")
		}
		// Test that Close is idempotent
		ch.Close()
	})
}
