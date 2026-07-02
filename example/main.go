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
