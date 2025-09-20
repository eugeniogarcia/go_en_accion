package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	usingMapas()
	/*declarations()
	usingSlices()
	arrayConversions()
	arrayPointerConversions()
	//panicArrayConversions()
	usingMapas()
	*/
}

func declarations() {
	fmt.Println("Usando declarations()\n")

	var a1 = [...]int{10, 20, 30} // con ... el compilador infiere el tamaño del array
	var a2 = [3]int{10, 20, 30}
	var b = [4]int{10, 20, 30} // con ... el compilador infiere el tamaño del array
	var c1 = []int{10, 20, 30} // declara un slice. La estructura del slice con el runtime de go es un struct con un puntero al array, la longitud y la capacidad
	var c2 = []int{10, 20, 30}

	fmt.Println(a1)       //[10 20 30]
	fmt.Println(a1 == a2) //Los arrays son comparables si son del mismo tipo y tamaño. Los slices no son comparables

	fmt.Println(b)                    //[10 20 30 0]
	fmt.Println(c1)                   //[10 20 30]
	fmt.Println(slices.Equal(c1, c2)) //Los slices no son comparables. Para comparar dos podemos usar la función slices.Equal del paquete slices
}

func usingSlices() {
	fmt.Println("\nusando usingSlices()\n")

	//Los slices se modelan en el runtime de Go como un struct con un puntero al array, la longitud y la capacidad.
	var a1 []int                      //declara un slices. Ahora mismo valor nil, longitud 0 y capacidad 0
	var a2 []int = []int{3, 4}        //declara un slice e inicializa con dos elementos
	fmt.Println(a1, len(a1), cap(a1)) // [] 0 0
	fmt.Println(a2, len(a2), cap(a2)) // [3 4] 2 2

	//crea un slice con make. El primer parámetro es el tipo, el segundo la longitud y el tercero la capacidad
	//si no se indica el tercer parámetro, la capacidad es igual a la longitud
	//si la longitud es 0, el slice apunta a un array de tamaño 0
	a1 = make([]int, 2, 4)

	//asigna a2 a1. En este momento a1 y a2 son dos slices diferentes (dos estructuras con su puntero, longitud y capacidad), que apuntan a la misma ubicación de memoria (array subyacente), tienen la misma longitud y capacidad
	a2 = a1
	for i := 0; i < len(a1); i++ {
		a1[i] = i + 1
	}
	fmt.Println(a1, len(a1), cap(a1)) // [1 2] 2 4
	fmt.Println(a2, len(a2), cap(a2)) // [1 2] 2 4

	//si añadimos un elemento a a1 aumenta su longitud y a2 no cambia. Como a1 tiene capacidad para incluir un elemento más, no se crea un nuevo array subyacente.
	a1 = append(a1, 3)
	fmt.Println(a1, len(a1), cap(a1)) // [1 2 3] 3 4
	fmt.Println(a2, len(a2), cap(a2)) // [1 2] 2 4

	//si añadimos dos elementos más a a1, al superar la capacidad, se crea un nuevo array subyacente y a2 no cambia (porque sigue apuntando al array original)
	a1 = append(a1, 4, 5)
	a1[0] = 100
	fmt.Println(a1, len(a1), cap(a1)) // [100 2 3 4 5] 5 8
	fmt.Println(a2, len(a2), cap(a2)) // [1 2] 2 4

	//podemos limpiar un slice. Mantiene el slice, pero cambia el valor de los elementos en el array subyacente a su valor cero
	clear(a2)                         //limpia a2
	fmt.Println(a2, len(a2), cap(a2)) // [0 0] 2 4

	//slices de slices
	x := make([]string, 0, 10) //creamos un slice de strings con longitud 0 y capacidad 6
	x = append(x, "a", "b", "c", "d", "e")
	y := x[:2:2]  //la capacidad y la longitud son iguales
	z := x[2:4:6] //creamos el slice apuntando al array subyacente de x, con longitud 2, y forzamos que la capacidad llegue hasta la posición 6 del array subyacente. Como hemos empezado en la 2, la capacidad es 6-2=4
	w := x[:4:6]
	fmt.Println(x, len(x), cap(x)) // [a b c d] 5 10
	fmt.Println(y, len(y), cap(y)) // [a b] 2 2
	fmt.Println(z, len(z), cap(z)) // [c d] 2 4
	fmt.Println(w, len(w), cap(w)) // [a b c d] 4 6

	//copia. Copia un slice en otro. Copia un número de posiciones igual a la longitud del slice más pequeño
	fuente := []int{1, 2, 3, 4}
	destino := make([]int, 2)
	eltosCopiados := copy(destino, fuente)
	fmt.Println(fuente)        // [1 2 3 4]
	fmt.Println(destino)       // [1 2]
	fmt.Println(eltosCopiados) // 2
}

func arrayConversions() {
	fmt.Println("\nusando arrayConversions()\n")

	xSlice := []int{1, 2, 3, 4} //declara un slice
	xArray := [4]int(xSlice)    //crea un array a partir del slice. Los dos apuntan a diferentes ubicaciones de memoria
	smallArray := [2]int(xSlice)
	xSlice[0] = 10
	smallArray[0] = 20
	fmt.Println(xSlice)     //[10 2 3 4]
	fmt.Println(xArray)     //[1 2 3 4]
	fmt.Println(smallArray) //[20 2]
}

func arrayPointerConversions() {
	fmt.Println("\nusando arrayPointerConversions()\n")

	xSlice := []int{1, 2, 3, 4}
	xArrayPointer := (*[3]int)(xSlice) //crea un puntero de arrays a un array a partir del slice
	xSlice[0] = 10
	xArrayPointer[1] = 20
	fmt.Println(xSlice)        //[10 20 3 4]
	fmt.Println(xArrayPointer) //&[10 20 3]
}

func panicArrayConversions() {
	fmt.Println("\nusando panicArrayConversions()\n")

	xSlice := []int{1, 2, 3, 4}
	panicArray := [5]int(xSlice) //no se puede crear un array de tamaño diferente al slice
	fmt.Println(panicArray)
}

func usingMapas() {
	fmt.Println("\nusando usingMapas()\n")

	m := map[string]int{
		"hello": 5,
		"world": 0,
	}
	//si buscamos con una key que no existe, devuelve el valor cero del tipo del valor, no da error
	valor, ok := m["hello"]
	fmt.Println(valor, ok)

	valor, ok = m["world"]
	fmt.Println(valor, ok)

	valor, ok = m["goodbye"]
	fmt.Println(valor, ok)

	valor = m["goodbye"]
	fmt.Println(valor)

	//podemos borrar una entrada
	delete(m, "hello")
	//eliminar todas las entradas del mapa
	fmt.Println(m, len(m))
	clear(m)
	fmt.Println(m, len(m)) //en un slice al hacer clear se mantenia el array subyacente, no se liberaba memoria, la capacidad no se cambiaba. Tambien se hacia el len a 0. Con el mapa se libera la memoria. Tener en cuenta que un mapa es una hasha table, es decir, que no es un espacio contiguo de memoria como un array o slice

	//los mapas, como los slices no son comparables
	a := map[string]int{
		"hello": 5,
		"world": 10,
	}
	b := map[string]int{
		"world": 10,
		"hello": 5,
	}
	fmt.Println(maps.Equal(a, b)) //true. Los mapas no son comparables, pero podemos usar la función maps.Equal del paquete maps para comparar dos mapas

}
