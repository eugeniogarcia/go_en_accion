package main

import "fmt"

func main() {
	in1 := make(chan int)
	in2 := make(chan int)
	go func() {
		for i := 10; i < 100; i += 10 {
			in1 <- i
		}
		close(in1)
	}()
	go func() {
		for i := 20; i >= 0; i-- {
			in2 <- i
		}
		close(in2)
	}()
	result := leerDosCanales(in1, in2)
	fmt.Println(result)
}

func leerDosCanales(in, in2 chan int) []int {
	var out []int
	// leemos de forma indefinida de ambos canales
	for count := 0; count < 2; {
		select {
		case v, ok := <-in:
			if !ok {
				in = nil // impedimos así que en la siguiente iteración este case vuelva a ser true
				count++  // un canal menos que leer
				continue
			}
			out = append(out, v)
		case v, ok := <-in2:
			if !ok {
				in2 = nil // impedimos así que en la siguiente iteración este case vuelva a ser true
				count++   // un canal menos que leer
				continue
			}
			out = append(out, v)
		}
	}
	return out
}
