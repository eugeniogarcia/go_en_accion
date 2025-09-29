package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		inGoroutine := 1
		ch1 <- inGoroutine //esta gorutina se bloquea hasta que alguien lea del canal ch1
		fromMain := <-ch2  // leemos
		fmt.Println("goroutine:", inGoroutine, fromMain)
	}()

	inMain := 2
	ch2 <- inMain          // se bloquea hata que alguien haya leido del canal ch2
	fromGoroutine := <-ch1 //leemos de ch1. Esto desbloquearía la gorutina, pero para llegar aquí la gorutina tiene que ejecutarse primero, y esta bloqueada
	fmt.Println("main:", inMain, fromGoroutine)
}
