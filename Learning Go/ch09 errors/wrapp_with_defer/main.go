package main

import "fmt"

//demuestra una forma eficiente de wrappear errores en una función con múltiples puntos de fallo
//usando un defer al final de la función para wrappear el error devuelto

func doThing1(val int) (int, error) {
	if val < 0 {
		return 0, fmt.Errorf("val no puede ser negativo")
	}
	return val * 2, nil

}

func doThing2(val int) (int, error) {
	if val%2 == 0 {
		return 0, fmt.Errorf("val no puede ser par")
	}
	return val * 2, nil

}

func doThing3(val1, val2 int) (string, error) {
	if val1+val2 > 100 {
		return "", fmt.Errorf("la suma de val1 y val2 no puede ser mayor que 100")
	}
	return fmt.Sprintf("val1: %d, val2: %d", val1, val2), nil
}

func DoSomeThings(val1 int, val2 int) (_ string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("en DoSomeThings: %w", err)
		}
	}()
	val3, err := doThing1(val1)
	if err != nil {
		return "", err
	}
	val4, err := doThing2(val2)
	if err != nil {
		return "", err
	}
	return doThing3(val3, val4)
}
