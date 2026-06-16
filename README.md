# go-async-await

A tiny Go library that provides **syntactic sugar** for goroutines and channels, mimicking
JavaScript’s `async/await` style.  
 ⚠️ **Note:** This library does not add any performance benefit compared to using goroutines +
channels directly. It’s purely a wrapper to reduce mental overhead when switching between Go and
JavaScript async models.

---

## ✨ Motivation

- Go’s native concurrency (`go f()` + `<-ch`) is simple and powerful, but the syntax differs from
  JavaScript’s `async/await`.
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
    // Launch async tasks
    a1 := Async(func() (string, error) { return worker(1, 3*time.Second) })
    a2 := Async(func() (string, error) { return worker(2, 1*time.Second) })
    a3 := Async(func() (string, error) { return worker(3, 2*time.Second) })

    // Await results (blocking until each finishes)
    msg, err := a1.Await()
    fmt.Println("Result:", msg, "Error:", err)

    msg, err = a2.Await()
    fmt.Println("Result:", msg, "Error:", err)

    msg, err = a3.Await()
    fmt.Println("Result:", msg, "Error:", err)
}

```

or you can find it at this
[example](https://github.com/hakimjazuli/go-async-await/blob/main/example/example.go)

---

## 📖 API

- `Async[T](<func()> 'T, error') \*Promise[T]`

  > - Launches a function in a goroutine and returns a Promise.

- `(\*Promise[T]) Await() (T, error)`
  > - Blocks until the goroutine completes. Returns (cached result, cached error) if called multiple
  >   times.

---

## ⚠️ Notes

- This is syntactic sugar only. Under the hood, it’s just goroutines + channels.

- No performance gain compared to idiomatic Go concurrency.

- Intended for developers who frequently switch between Go and JS async code and want a consistent
  mental model.

---

## 📝 License

MIT
