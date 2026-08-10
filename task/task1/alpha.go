package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

// Функция сохраняет снимок профиля в файл
func profile(name string) {
	file, _ := os.Create("a" + ".prof")
	defer file.Close()
	p := pprof.Lookup(name)
	if p != nil {
		p.WriteTo(file, 0)
	}
}

func main() {
	// Включаем сбор профилей блокировок и мьютексов
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	const nGoro = 10
	const nOps = 100000

	var total int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for range nGoro {
		wg.Add(1) // Увеличиваем счетчик WaitGroup
		go func() {
			defer wg.Done()
			for range nOps {
				mu.Lock()
				total++
				mu.Unlock()
			}
		}()
	}

	wg.Wait() // Главный поток заблокируется здесь, ожидая горутины
	fmt.Println("total =", total)

	// Снимаем профиль строго ПОСЛЕ завершения вычислений
	profile("block")
}