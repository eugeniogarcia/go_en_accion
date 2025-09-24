package main

import (
	"errors"
	"fmt"
)

// define un tipo que implementa error. El error no es comparable, no podremos usar errors.Is, salvo que implementemos el métod IS, pero si podemos usar errors.As
type MyErr struct {
	Codes []int
}

func (me MyErr) CodeVals() []int {
	return me.Codes
}

func (me MyErr) Error() string {
	return fmt.Sprintf("codigos: %v", me.Codes) //usamos %v
}

func creaErrorMyErr() error {
	return MyErr{Codes: []int{1, 1, 2, 3, 5, 8}}
}

func main() {
	err := creaErrorMyErr()
	var miError MyErr

	// para utilizar AS hay que pasar un puntero al tipo que queremos comprobar, independientemete de que tenga un valor asignado o no
	if errors.As(err, &miError) {
		fmt.Println(miError.Codes)
	}

	//podemos usar AS con cualquier tipo o interfaz. Aquí comprobamos si err implementa la interfaz anonima que define el método CodeVals
	var coder interface {
		CodeVals() []int
	}
	if errors.As(err, &coder) {
		fmt.Println(coder.CodeVals())
	}

}
