package goasyncawait

import (
	"fmt"
)

// Result is a generic container that holds either:
// - a successful value of type T
// - or an error if the computation failed
type Result[T any] struct {
	val T     // the actual value returned
	err error // any error encountered during execution
}

// Promise represents an asynchronous computation.
// It wraps a channel that delivers a Result[T] once the goroutine finishes.
// It also caches the result so multiple Await calls don’t re-block.
type Promise[T any] struct {
	ch     <-chan Result[T] // channel carrying the result from the goroutine
	called bool             // tracks whether Await has already been called
	result Result[T]        // cached result after first Await
}

// Await blocks until the goroutine sends a value into the channel.
// On first call, it consumes from the channel and caches the result.
// On subsequent calls, it returns the cached result immediately.
func (p *Promise[T]) Await() (T, error) {
	if p.called {
		// If Await was already called, return cached result
		return p.result.val, p.result.err
	}
	// Block until a value is received from the channel
	p.result = <-p.ch
	p.called = true
	return p.result.val, p.result.err
}

// Then registers a callback to be executed once the current Promise resolves.
// It executes the callback in a non-blocking, asynchronous manner and returns
// a new Promise containing the transformed result.
func (p *Promise[T]) Then(callback func(T, error) (T, error)) *Promise[T] {
	// We use Async to immediately return a new Promise without blocking the caller.
	return Async(func() (T, error) {
		// Await the current promise's resolution (this happens inside the background goroutine)
		val, err := p.Await()

		// Pass the results to your callback function and return its outcomes
		return callback(val, err)
	})
}

// Async launches the provided function in a new goroutine.
// The function must return (T, error).
// Async captures panics and converts them into an error Result.
// It returns a Promise[T] that can be awaited later.
func Async[T any](f func() (T, error)) *Promise[T] {
	ch := make(chan Result[T], 1) // buffered channel to hold one result
	go func() {
		// Recover from panic and send it as an error
		defer func() {
			if rec := recover(); rec != nil {
				var zero T
				ch <- Result[T]{val: zero, err: fmt.Errorf("panic: %v", rec)}
			}
		}()
		// Execute the function and send its result
		val, err := f()
		ch <- Result[T]{val: val, err: err}
	}()
	// Return a Promise wrapping the channel
	return &Promise[T]{ch: ch}
}
