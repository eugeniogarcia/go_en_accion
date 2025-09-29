package main

import (
	"context"
	"fmt"
)

func countTo(ctx context.Context, max int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < max; i++ {
			select {
			case <-ctx.Done(): //2. este método devuelve un canal que se cierra cuando el contexto se cancela
				return
			case ch <- i:
			}
		}
	}()
	return ch
}

func main() {
	ctx, cancel := context.WithCancel(context.Background()) //1. creamos un contexto cancelable. Devuelve el contexto y una función para cancelarlo

	//pasaremos el contexto a la gorutina
	ch := countTo(ctx, 10)
	for i := range ch {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}
	cancel() //3. llamamos a la función para cancelar el contexto
}
