package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	result, err := controlaTiempo(haceTrabajo, 2*time.Second)
	fmt.Println(result, err)
}

func controlaTiempo[T any](worker func() T, limit time.Duration) (T, error) {
	out := make(chan T, 1)
	//contexto que implementa un timeout
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	// aseguramos que el contexto se cancele al finalizar la función
	defer cancel()

	//lanzamos la gorutina
	go func() {
		out <- worker()
	}()

	//observamos el canal de salida y el contexto
	select {
	case result := <-out:
		return result, nil
	case <-ctx.Done(): //habrá un token cuando se alcance el timeout o se cancele el contexto
		var zero T
		return zero, errors.New("timed out")
	}
}

func haceTrabajo() int {
	if x := rand.Int(); x%2 == 0 {
		return x
	} else {
		time.Sleep(10 * time.Second)
		return 100
	}
}
