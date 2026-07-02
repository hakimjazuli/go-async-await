# go-async-await

A tiny Go library that provides **syntactic sugar** for goroutines and channels, mimicking
JavaScript’s `async/await/then` style.  
 ⚠️ **Note:** This library does not add any performance benefit compared to using goroutines +
channels directly. It’s purely a wrapper to reduce mental overhead when switching between Go and
JavaScript async models.

---

## ✨ Motivation

- Go’s native concurrency (`go f()` + `<-ch`) is simple and powerful, but the syntax differs from
  JavaScript’s `async/await/then`.
- For JS developers working in mixed stacks, constantly switching mental models can be tiring:
  - Fullstack projects (Go backend + JS frontend)
  - Desktop apps using [Wails](https://wails.io/) (Go native binding IPC + JS UI)
- This library provides a **Promise-like API** in Go, so you can write async code with a familiar
  flow.

---

## 📦 Installation

```bash
go get github.com/hakimjazuli/go-async-await/goasyncawait
```

---

## 🧭 example

```go
package main

import (
	"fmt"
	"time"

	. "github.com/hakimjazuli/go-async-await/goasyncawait"
)

func worker(id int, sleep time.Duration) (string, error) {
	fmt.Println("before sleep, id:", id, "sleep:", sleep)
	time.Sleep(sleep)
	fmt.Println("after sleep, id:", id, "sleep:", sleep)
	if id == 2 {
		return "", fmt.Errorf("worker %d failed", id)
	}
	return fmt.Sprintf("Worker %d finished after %v", id, sleep), nil
}

func main() {
	// 1. Fire off the original worker promises
	a1 := Async(func() (string, error) { return worker(1, 3*time.Second) })
	a2 := Async(func() (string, error) { return worker(2, 1*time.Second) })
	a3 := Async(func() (string, error) { return worker(3, 2*time.Second) })

	// 2. Chain asynchronous operations using .Then()
	// These will run automatically in the background as soon as their worker resolves.
	// Setting these up is completely non-blocking.
	chainedA1 := a1.Then(func(msg string, err error) (string, error) {
		if err != nil {
			return "", fmt.Errorf("chained1 caught error: %w", err)
		}
		return msg + " -> [Chained Modification 1]", nil
	})

	chainedA2 := a2.Then(func(msg string, err error) (string, error) {
		if err != nil {
			// Worker 2 fails, so this error handling block will execute
			return "Fallback value for failed Worker 2", nil
		}
		return msg + " -> [Chained Modification 2]", nil
	})

	chainedA3 := a3.Then(func(msg string, err error) (string, error) {
		if err != nil {
			return "", fmt.Errorf("chained3 caught error: %w", err)
		}
		return msg + " -> [Chained Modification 3]", nil
	})

	fmt.Println("--- Main thread continues setup uninterrupted ---")

	// 3. Await the final chained results.
	// Only these lines will block execution until the background work completes.
	msg, err := chainedA1.Await()
	fmt.Println("Final Chained 1 Result:", msg, "| Error:", err)

	msg, err = chainedA2.Await()
	fmt.Println("Final Chained 2 Result:", msg, "| Error:", err)

	msg, err = chainedA3.Await()
	fmt.Println("Final Chained 3 Result:", msg, "| Error:", err)
}

```

or you can find it at this
[example](https://github.com/hakimjazuli/go-async-await/blob/main/example/main.go)

---

## 📖 API

- `Async[T any](f func() (T, error)) *Promise[T]`

  > - Launches the provided function in a new background goroutine.
  > - Automatically catches panics and converts them into errors.
  > - Returns a `*Promise[T]` that can be awaited or chained.

- `(*Promise[T]) Await() (T, error)`

  > - Blocks the current goroutine until the asynchronous operation completes.
  > - Caches the results; subsequent calls return the cached value and error immediately without
  >   re-blocking.

- `(*Promise[T]) Then(callback func(T, error) (T, error)) *Promise[T]`
  > - Registers a callback to process the result or handle the error of the current promise.
  > - **Non-blocking:** Executes asynchronously in a new background goroutine, allowing immediate
  >   chaining (`.Then().Then()`).
  > - Returns a new `*Promise[T]` carrying the transformed outcome.

---

## ⚠️ Notes

- This is syntactic sugar only. Under the hood, it’s just goroutines + channels.

- No performance gain compared to idiomatic Go concurrency.

- Intended for developers who frequently switch between Go and JS async code and want a consistent
  mental model.

---

## 📝 License

MIT

## 📌 Version

Current stable release: **v0.2.0**

Install with:

```bash
go get github.com/hakimjazuli/go-async-await/goasyncawait@v0.2.0
```
