package main

import (
	"net/http"
	"time"
)

type HolaMundoHandler struct{}

// Hacemos que HolaMundoHandler implemente la interface http.Handler
func (hh HolaMundoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola!\n")) //enviamos la respuesta
}

func main() {
	s := http.Server{
		Addr:         ":8080", //host y puerto en el que escuchamos
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      HolaMundoHandler{}, //handler. La estructura HolaMundoHandler implementa la interface Handler
	}
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
