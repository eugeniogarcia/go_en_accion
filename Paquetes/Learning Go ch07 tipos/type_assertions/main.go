package main

import "fmt"

// Definimos un nuevo tipo basado en int
type MyInt int

func main() {
	typeAssert()
	typeAssertPanicTypeNotIdentical()

	err := typeAssertCommaOK()
	if err != nil {
		fmt.Println(err)
	}
}

// Ejemplo de type assertion usando la forma comma-ok para evitar pánico
func typeAssertCommaOK() error {
	var i any           // i es una interfaz que aplica a cualquier tipo
	var mine MyInt = 20 //mine es de tipo MyInt
	i = mine            // Asignamos mine (MyInt) a i (any)
	i2, ok := i.(int)   // Intentamos afirmar que i es un int. En lugar de crear pánico, ok será false si no es del tipo esperado
	if !ok {
		return fmt.Errorf("unexpected type for %v", i) // int y MyInt no son idénticos, devolvemos un error
	}
	fmt.Println(i2 + 1)
	return nil
}

// Ejemplo que provoca pánico por intentar afirmar tipos no idénticos
func typeAssertPanicTypeNotIdentical() {
	defer func() {
		if m := recover(); m != nil {
			fmt.Println(m) // Se imprime porque ocurre un pánico
		}
	}()

	var i any // i es una interfaz que aplica a cualquier tipo
	var mine MyInt = 20
	i = mine      // Asignamos mine (MyInt) a i (any)
	i2 := i.(int) // Esto provoca pánico porque MyInt no es idéntico a int
	fmt.Println(i2 + 1)
}

// Ejemplo exitoso de type assertion
func typeAssert() {
	var i any // i es una interfaz que aplica a cualquier tipo
	var mine MyInt = 20
	i = mine        // Asignamos mine (MyInt) a i (any)
	i2 := i.(MyInt) // Afirmamos correctamente que i es MyInt
	fmt.Println(i2 + 1)
}
