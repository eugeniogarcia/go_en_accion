package main

import "fmt"

func main() {
	deferExample()
}

//vamos a etiquetar la variable de respuesta. Esto declara una variable local que podemos informar con el valor a retornar. En todo caso si no la usamos el runtime igual su valor al valor incluido en return
func deferExample() (respuesta int) {
	a := 10
	defer func(val int) {
		fmt.Println("first:", val)
		fmt.Println("respuesta:", respuesta)
	}(a)
	a = 20
	defer func(val int) {
		fmt.Println("second:", val)
	}(a)
	a = 30
	fmt.Println("exiting:", a)
	return a
}
