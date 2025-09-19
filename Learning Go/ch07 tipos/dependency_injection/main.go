package main

import (
	"fmt"
	"net/http"
)

// Tipo que es cumple con la interface Logger
type LoggerAdapter func(message string)

func (lg LoggerAdapter) Log(message string) {
	lg(message)
}

func LogOutput(message string) {
	fmt.Println(message)
}

func main() {
	// Creamos una instancia de LoggerAdapter. Cumple con la interface Logger
	l := LoggerAdapter(LogOutput)

	//Cremos una instancia de SimpleDataStore. Cumple con la interface DataStore
	var ds SimpleDataStore = NewSimpleDataStore()

	//Cremos una instancia de SimpleLogic. La factoria espera dos interfaces Logger y DataStore. Cumple con la interface Logic
	var logic SimpleLogic = NewSimpleLogic(l, ds)

	//Creamos una instancia de Controller, que cumple con la interface Controller. La factoria espera dos interfaces Logger y Logic
	var c Controller = NewController(l, logic)

	//Creamos un servidor http, y usamos como handler nuestra lógica de negocio
	http.HandleFunc("/hola", c.SayHello)
	http.ListenAndServe(":8080", nil)
}
