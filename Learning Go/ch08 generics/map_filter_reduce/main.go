package main

import (
	"fmt"
)

// declaramos una funcion Map con dos parametros de cualquier tipo. Tratara un slice de T1, y aplicara una función para devolver un slice de tipo T2
func Mapa[T1, T2 any](s []T1, f func(T1) T2) []T2 {
	r := make([]T2, len(s)) //creamos el slice de salida de tipo T2
	for i, v := range s {
		r[i] = f(v) //aplicamos la función f sobre cada elemento del slice de entrada para obtener el elemento del slice de salida
	}
	return r
}

// Reductor que se implementa con dos tipos genéricos T1 y T2. Toma un slice de T1, un valor inicial de T2 y una función que toma un acumulador de T2 y un valor de T1 y devuelve un nuevo acumulador de T2.
func Reductor[T1, T2 any](s []T1, initializer T2, f func(T2, T1) T2) T2 {
	r := initializer // el valor inicial, que es de tipo T2
	for _, v := range s {
		r = f(r, v) //cada elemento del slice T1 se trata con la función y se actualiza el acumulador r
	}
	return r
}

// Filter filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
// Aplica un filtr sobre el slice de tipo T. El filtro consiste en una función que toma un T y devuelve un booleano. Si la función devuelve true, el elemento se incluye en el slice de salida.
func Filtro[T any](s []T, f func(T) bool) []T {
	var r []T //creamos el slice de salida
	for _, v := range s {
		if f(v) {
			r = append(r, v) // añadimos el elemento al slice de salida
		}
	}
	return r
}

func main() {
	words := []string{"One", "Potato", "Two", "Potato"}
	filtered := Filtro(words, func(s string) bool {
		return s != "Potato"
	})
	fmt.Println(filtered)
	lengths := Mapa(filtered, func(s string) int {
		return len(s)
	})
	fmt.Println(lengths)
	sum := Reductor(lengths, 0, func(acc int, val int) int {
		return acc + val
	})
	fmt.Println(sum)
}
