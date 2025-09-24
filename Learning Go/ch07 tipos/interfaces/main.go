package main

import (
	"errors"
	"fmt"
)

type myinterface interface {
	f(string) string
}

type mystruct struct {
	men string
}

func (a mystruct) f(s string) string {
	return fmt.Sprintf("%s %s", s, a.men)
}

type mystruct2 struct {
	men string
}

func (a *mystruct2) f(s string) string {
	return fmt.Sprintf("%s %s", s, a.men)
}

func pruebaInterfaces(caso int) (myinterface, error) {
	var x myinterface // x es una interfaz nil
	var y mystruct    // y es una estructura nil
	switch caso {
	case 0:
		return x, nil //devolvemos una interfaz nil
	case 1:
		return y, nil // devolvemos una interfaz no nil que contiene una estructura nil. Recordemos que una interfaz es un par (tipo, valor). En este caso el tipo es mystruct y el valor es nil
	case 2:
		return mystruct{men: "Hola"}, nil //mystruct implementa myinterface
	case 3:
		return &mystruct{men: "Hola"}, nil //&mystruct es un puntero a a mystruct que implementa myinterface
	case 4:
		return &mystruct2{men: "Hola"}, nil //&mystruct2 es un puntero a a mystruct2 que implementa myinterface. Notese que mystruct2 no implementa myinterface (porque el receptor es un puntero)
	default:
		return x, errors.New("caso no soportado")
	}
}

func main() {
	for i := 0; i < 6; i++ {
		v, e := pruebaInterfaces(i)
		if e != nil {
			fmt.Println(e)
		} else {
			switch v := v.(type) { //como v es una interfaz, podemos hacer type assertion
			case nil: //una interfaz nil
				fmt.Println("valor nulo")
			case mystruct: //una interfaz que hace referencia al tipo mystruct
				fmt.Println("mystruct")
				fmt.Println(v.f("Eugenio"))
			case *mystruct: //una interfaz que hace referencia a un puntero a mystruct
				fmt.Println("*mystruct")
				fmt.Println(v.f("Eugenio"))
			case *mystruct2: //una interfaz que hace referencia a un puntero a mystruct2
				fmt.Println("*mystruct2")
				fmt.Println(v.f("Eugenio"))
			}
		}
	}
}
