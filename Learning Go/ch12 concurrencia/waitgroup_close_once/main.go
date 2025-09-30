package main

import (
	"fmt"
	"sync"
)

func procesaVarios[T, R any](in <-chan T, processor func(T) R, num int) []R {
	out := make(chan R, num)
	var wg sync.WaitGroup //creamos un waitgroup

	wg.Add(num) //incrementamos el contador del waitgroup

	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done() //decrementamos el contador del waitgroup
			for v := range in {
				out <- processor(v)
			}
		}()
	}
	go func() {
		wg.Wait() //esperamos a que el contador del waitgroup llegue a 0
		close(out)
	}()
	var result []R
	for v := range out {
		result = append(result, v)
	}
	return result
}

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 20; i++ {
			ch <- i
		}
		close(ch)
	}()
	results := procesaVarios(ch, func(i int) int {
		return i * 2
	}, 3)
	fmt.Println(results)
}
