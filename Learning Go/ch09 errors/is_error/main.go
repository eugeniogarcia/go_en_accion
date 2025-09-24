package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
)

// define un tipo que implementa error. Como el tipo incluye un slice, y el slice no es comparable, no se puede usar == para comparar errores de este tipo, esto es, errors.Is no va a funcionar salvo que en la definicion del tipo se implemente el metodo Is
type MyErr struct {
	Codes []int
}

func (me MyErr) Error() string {
	return fmt.Sprintf("codes: %v", me.Codes)
}

// implementamos el método Is para que errors.Is pueda comparar errores de este tipo
func (me MyErr) Is(target error) bool {
	if me2, ok := target.(MyErr); ok { //comprobamos si error tiene como underlying type MyErr, y en caso afirmativo aplicamos la lógica de comparación
		return slices.Equal(me.Codes, me2.Codes) //en este caso usamos slices.Equal para comparar los slices
	}
	return false
}

func creaErrorMyErr() error {
	return fmt.Errorf("error MyErr empaquetado: %w", MyErr{
		Codes: []int{2, 7, 1, 8, 2, 8},
	})
}

func creaSentinelError() error {
	f, err := os.Open("not_existe.txt")
	if err != nil {
		return fmt.Errorf("en creaSentinelError: %w", err)
	}
	f.Close()
	return nil
}

func main() {
	err := creaSentinelError()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("El archivo no existe")
		}
	}

	err = creaErrorMyErr()
	miError1 := MyErr{Codes: []int{2, 7, 1, 8, 2, 8}}
	if errors.Is(err, miError1) {
		fmt.Println("Lo encontre!")
	}

	miError2 := MyErr{Codes: []int{2, 7, 1, 8, 2, 8}}
	if errors.Is(err, miError2) {
		fmt.Println("Lo encontre 2!")
	} else {
		fmt.Println("No es el error que busco")
	}
}
