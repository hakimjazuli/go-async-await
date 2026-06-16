package main

import (
	"fmt"
	"time"

	"github.com/hakimjazuli/go-async-await/goasyncawait"
)

func worker(id int, sleep time.Duration) (string, error) {
	time.Sleep(sleep)
	if id == 2 {
		return "", fmt.Errorf("worker %d failed", id)
	}
	return fmt.Sprintf("Worker %d finished after %v", id, sleep), nil
}

func main() {
	p := goasyncawait.Async(func() (string, error) { return worker(1, 2*time.Second) })
	msg, err := p.Await()
	fmt.Println("Result:", msg, "Error:", err)
}
