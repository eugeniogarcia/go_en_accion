package main

import (
	"fmt"
)

// Declaramos una interfaz con un método A que recibe un int y devuelve un int
type miInterface interface {
	A(int) int
}

// Declaramos un tipo función, que implementa la interfaz miInterface
type miFun func(int) int

func (f miFun) A(i int) int {
	return f(i)
}

// Declaramos dos funciones simples
func doble(v int) int {
	return v * 2
}

func triple(v int) int {
	return v * 3
}

// Función que recibe una función como parámetro
func paraguas(f miFun, v int) int {
	return f(v)
}

// Función que recibe una interfaz como parámetro
func paraguas2(f miInterface, v int) int {
	return f.A(v)
}

func main() {
	// Pasamos la función como argumento a paraguas (usa el tipo miFun)
	fmt.Println(paraguas(doble, 2))  // Imprime 4
	fmt.Println(paraguas(triple, 2)) // Imprime 6

	// No podemos pasar las funciones directamente a paraguas2 porque no implementan la interfaz
	// fmt.Println(paraguas2(doble, 2))
	// fmt.Println(paraguas2(triple, 2))

	// Pero sí podemos convertirlas a miFun, que implementa la interfaz miInterface
	var f miFun = doble
	fmt.Println(paraguas2(f, 2)) // Imprime 4

	// O convertir directamente en la llamada
	fmt.Println(paraguas2(miFun(triple), 2)) // Imprime 6

	// La función no implementa la interfaz, así que esto no compila:
	// var f1 interfaz = doble

	// Pero sí podemos convertirla a miFun, que implementa la interfaz
	var f1 miInterface = miFun(doble)

	// Podemos hacer un type assertion para comprobar que f1 es de tipo miFun
	_, ok := f1.(miFun)
	fmt.Println(ok) // Imprime true
}
