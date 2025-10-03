package main

import (
	"net/http"
	"time"
)

/*
Endpoints:
/person/greet
/dog/greet
/hello
*/
func main() {
	persona := http.NewServeMux() //Creamos un multiplexor para el subruta /saludos
	persona.HandleFunc("/saludos", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("saludos!\n"))
	})

	perro := http.NewServeMux() //Creamos un multiplexor para el subruta /saludos
	perro.HandleFunc("/saludos", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("buen chucho!\n"))
	})

	multiplexador := http.NewServeMux() //Creamos un multiplexor...
	//http.StripPrefix elimina el prefijo de la ruta antes de pasarla al multiplexor asignado, es decir, si llamamos a /persona/saludos, quitamos /persona, de modo que el multiplexor "persona" recibirá /saludos
	multiplexador.Handle("/persona/", http.StripPrefix("/persona", persona))         //... para la ruta /persona/. Le asignamos el multplexor persona
	multiplexador.Handle("/perro/", http.StripPrefix("/perro", perro))               //... para la ruta /perror/. Le asignamos el multplexor perro
	multiplexador.HandleFunc("/hola", func(w http.ResponseWriter, r *http.Request) { // ... para la ruta /hola. Le asignamos un handler directamente
		w.Write([]byte("Hola Mundo!\n"))
	})

	/*
		/persona/saludos -> "saludos!"
		/perro/saludos -> "buen chucho!"
		/hola -> "Hola Mundo!"
	*/
	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      multiplexador,
	}
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
