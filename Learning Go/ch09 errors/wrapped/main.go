package main

import (
	"errors"
	"fmt"
	"os"
)

type Status int

// Error a medida que ademas permite Wrappear un error
type StatusErr struct {
	Status  Status
	Message string
	Err     error //para guardar el error wrappeado
}

// para implementar la interface error (tanto el tipo como el puntero a ese tipo implementan la interface error)
func (se StatusErr) Error() string {
	return se.Message
}

// para hacer unwrap del error wrappeado
func (se StatusErr) Unwrap() error {
	return se.Err
}

const (
	InvalidLogin Status = iota + 1
	NotFound
)

// Error a medida que ademas que permite Wrappear más de un error
type MyError struct {
	Code   int
	Errors []error //slice de errores para wrappear varios errores
}

// para implementar la interface error (tanto el tipo como el puntero a ese tipo implementan la interface error)
func (m MyError) Error() string {
	return errors.Join(m.Errors...).Error() //usamos errors.Join para unir los mensajes de los errores wrappeados
}

// para hacer unwrap de los errores wrappeados (usamos un slice de error)
func (m MyError) Unwrap() []error {
	return m.Errors
}

// Función para lanzar un error a medida
// lanzamos un error a medida que wrapea un error
func funcionWrapperror() error {
	return StatusErr{
		Status:  NotFound,
		Message: "fichero no encontrado",
		Err:     errors.New("error empaquetado"),
	}
}

// Función para lanzar un error a medida
// lanzamos un error a medida que wrapea dos errores
func funcionWrappMultierror() error {
	return MyError{
		Code: 12,
		Errors: []error{
			StatusErr{
				Status:  NotFound,
				Message: "fichero no encontrado",
			},
			errors.New("error empaquetado"),
		},
	}
}

// Creamos un error estandar usando %w para wrappear otro error
func fileChecker(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("en fileChecker: %w", err)
	}
	f.Close()
	return nil
}

// helper para tratar errores
func tratarError(err error) {
	switch err := err.(type) {
	case interface{ Unwrap() error }:
		// solo hay un error wrapeado. Lo extraemos
		innerErr := err.Unwrap()
		fmt.Println(innerErr)
	case interface{ Unwrap() []error }:
		// hay unvarios errores wrapeados. Los extraemos
		innerErrs := err.Unwrap()
		for _, innerErr := range innerErrs {
			fmt.Println(innerErr)
		}
	default:
		fmt.Println(err)
	}
}

func main() {
	var err error

	//demos
	err = fileChecker("not_here.txt")
	if err != nil {
		fmt.Println(err)
		if wrappedErr := errors.Unwrap(err); wrappedErr != nil {
			fmt.Println(wrappedErr)
		}
	}

	err = funcionWrappMultierror()
	tratarError(err)
	err = funcionWrapperror()
	tratarError(err)
}
