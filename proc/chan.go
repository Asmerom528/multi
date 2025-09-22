package proc

import (
	"encoding/gob"
	"os"
)

// Chan provides a type-safe channel for cross-process communication.
// Uses Unix pipes for IPC and gob for serialization.
// Not safe for concurrent use by multiple goroutines.
type Chan[T any] struct {
	w      *os.File
	r      *os.File
	enc    *gob.Encoder
	dec    *gob.Decoder
	closed bool
}

// NewChan creates a new Chan for cross-process communication.
func NewChan[T any]() *Chan[T] {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	return &Chan[T]{
		enc: gob.NewEncoder(w),
		dec: gob.NewDecoder(r),
		w:   w,
		r:   r,
	}
}

// Send encodes and sends a value of type T through the channel.
// Blocks if the pipe buffer is full.
func (c *Chan[T]) Send(val T) {
	err := c.enc.Encode(val)
	if err != nil {
		panic(err)
	}
}

// Recv receives and decodes a value of type T from the channel.
// Blocks if no data is available.
func (c *Chan[T]) Recv() T {
	var val T
	if c.closed {
		return val
	}
	err := c.dec.Decode(&val)
	if err != nil {
		panic(err)
	}
	return val
}

// Close closes the channel's file descriptors.
func (c *Chan[T]) Close() {
	if c.w != nil {
		_ = c.w.Close()
		c.w = nil
	}
	if c.r != nil {
		_ = c.r.Close()
		c.r = nil
	}
	c.closed = true
}

// Closed returns true if the channel has been closed.
func (c *Chan[T]) Closed() bool {
	return c.closed
}
