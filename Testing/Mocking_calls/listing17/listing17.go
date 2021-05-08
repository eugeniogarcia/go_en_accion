// This sample code implement a simple web service.
package main

import (
	"log"
	"net/http"

	"www.ejemplo.com/ejemplo/Mocking_calls/listing17/handlers"
)

// main is the entry point for the application.
func main() {
	handlers.Routes()

	log.Println("listener : Started : Listening on :4000")
	// Empieza a escuchar los end-points que hemos mockeado
	http.ListenAndServe(":4000", nil)
}
