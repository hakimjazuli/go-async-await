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
	a1 := Async(func() (string, error) { return worker(1, 3*time.Second) })
	a2 := Async(func() (string, error) { return worker(2, 1*time.Second) })
	a3 := Async(func() (string, error) { return worker(3, 2*time.Second) })

	msg, err := a1.Await()
	fmt.Println("Result:", msg, "Error:", err)

	msg, err = a2.Await()
	fmt.Println("Result:", msg, "Error:", err)

	msg, err = a3.Await()
	fmt.Println("Result:", msg, "Error:", err)
}
