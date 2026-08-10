package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"sync/atomic"
)

// Функция сохраняет снимок профиля в файл
func profile(name string) {
	file, _ := os.Create("b" + ".prof")
	defer file.Close()
	p := pprof.Lookup(name)
	if p != nil {
		p.WriteTo(file, 0)
	}
}

func main() {
	// Включаем сбор профилей блокировок
	runtime.SetBlockProfileRate(1)

	const nGoro = 10
	const nOps = 100000

	var total atomic.Int64
	var wg sync.WaitGroup

	for range nGoro {
		wg.Add(1) // Увеличиваем счетчик WaitGroup
		go func() {
			defer wg.Done()
			for range nOps {
				total.Add(1) // Атомарная операция (не блокирует горутину)
			}
		}()
	}

	wg.Wait() // Главный поток заблокируется здесь, но лишь на мгновение
	fmt.Println("total =", total.Load())

	// Снимаем профиль строго ПОСЛЕ завершения вычислений
	profile("block")
}