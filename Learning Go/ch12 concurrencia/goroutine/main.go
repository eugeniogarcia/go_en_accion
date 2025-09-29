package main

import "fmt"

func main() {
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	result := processConcurrently(x)
	fmt.Println(result)
}

func process(val int) int {
	// do something with val
	return val * 2
}

const numGoroutines = 5

func processConcurrently(inVals []int) []int {
	// creamos canales con un buffer igual al número de gorutinas
	in := make(chan int, numGoroutines)
	out := make(chan int, numGoroutines)
	// creamos y lanzamos las gorutinas
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for val := range in {
				result := process(val)
				out <- result
			}
		}() //notese que estamos haciendo una llamada
	}

	//han quedado lanzadas las gorutinas, ahora vamos a informar valores en el canal
	go func() {
		for _, v := range inVals {
			in <- v //se bloquea si el canal estuviera lleno
		}
		close(in) // cerramos el canal, ya no vamos a escribir más
	}()
	// creamos un slice para recoger los resultados
	outVals := make([]int, 0, len(inVals))
	for i := 0; i < len(inVals); i++ {
		//añadimos la respuesta del canal al slice; Se bloquea si el canal está vacío
		outVals = append(outVals, <-out)
	}
	return outVals
}
