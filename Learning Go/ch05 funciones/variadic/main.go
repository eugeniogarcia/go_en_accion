package main

import "fmt"

//indicamos un varadic con ...tipo. En el cuerpo de la fancion el argumento llega como un slice que captura todos los valores pasados. Podemos desempaquetar un slice con ...slice
func addTo(base int, vals ...int) []int {
	//creamos un slice con len 0 y capacidad igual al número de valores pasados
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		out = append(out, base+v)
	}
	return out
}

func main() {
	fmt.Println(addTo(3))
	fmt.Println(addTo(3, 2))
	fmt.Println(addTo(3, 2, 4, 6, 8))
	a := []int{4, 3}
	fmt.Println(addTo(3, a...))
	fmt.Println(addTo(3, []int{1, 2, 3, 4, 5}...))
}
