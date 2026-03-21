// Ограничитель скорости
package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrCanceled error = errors.New("canceled")

// начало решения

// throttle следит, чтобы функция fn выполнялась не более limit раз в секунду.
func throttle(limit int, fn func()) (handle func() error, cancel func()) {
	interval := time.Second / time.Duration(limit)
	ticker := time.NewTicker(interval)

	done := make(chan struct{})

	stopOnce := make(chan struct{}, 1)

	go func() {
		<-stopOnce
		ticker.Stop()
		close(done)
	}()

	handle = func() error {
		select {
		case <-done:
			return ErrCanceled
		case <-ticker.C:
			go fn()
			return nil
		}
	}

	cancel = func() {
		select {
		case stopOnce <- struct{}{}:
		default:
		}
	}
	return handle, cancel
}

// конец решения
func main() {
	work := func() {
		fmt.Print(".")
	}

	handle, cancel := throttle(5, work)
	defer cancel()

	start := time.Now()
	const n = 10
	for i := 0; i < n; i++ {
		handle()
	}
	fmt.Println()
	fmt.Printf("%d queries took %v\n", n, time.Since(start))
}
