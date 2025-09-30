package main

import (
	"fmt"
	"sync"
)

// declaramos a nivel de paquete las variables, para asegurar que todas las funciones hagan referencia a las mismas instancias
var parser ParserMuyComplejo
var once sync.Once

func Parse(dataToParse string) string {
	//aseguramos que esta función se ejecute una sola vez
	once.Do(func() {
		parser = initParser()
	})
	return parser.Parse(dataToParse)
}

type ParserMuyComplejo interface {
	Parse(string) string
}

func initParser() ParserMuyComplejo {
	// do all sorts of setup and loading here
	fmt.Println("initializing!")
	return SCPI{}
}

type SCPI struct {
}

func (s SCPI) Parse(in string) string {
	if len(in) > 1 {
		return in[0:1]
	}
	return ""
}

func main() {
	// "initializing!" will print out only once
	result := Parse("hello")
	fmt.Println(result)
	result2 := Parse("goodbye")
	fmt.Println(result2)
}
