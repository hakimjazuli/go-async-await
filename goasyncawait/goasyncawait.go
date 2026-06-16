package goasyncawait

import (
	"fmt"
)

// Result holds either a value or an error
type Result[T any] struct {
	val T
	err error
}

// Promise wraps a channel and exposes .Await()
type Promise[T any] struct {
	ch     <-chan Result[T]
	called bool
	result Result[T]
}

// await blocks until the goroutine sends a value.
// If called more than once, returns the cached result/error.
func (p *Promise[T]) Await() (T, error) {
	if p.called {
		return p.result.val, p.result.err
	}
	p.result = <-p.ch
	p.called = true
	return p.result.val, p.result.err
}

// Async takes a thunk (zero‑arg function) so the work is deferred
// until the goroutine runs, preserving concurrency.
func Async[T any](f func() (T, error)) *Promise[T] {
	ch := make(chan Result[T], 1)
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				var zero T
				ch <- Result[T]{val: zero, err: fmt.Errorf("panic: %v", rec)}
			}
		}()
		val, err := f()
		ch <- Result[T]{val: val, err: err}
	}()
	return &Promise[T]{ch: ch}
}
