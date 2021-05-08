// Sample test to show how to write a basic example.
package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

// ExampleSendJSON provides a basic example.
func ExampleSendJSON() {
	//Vamos a probar el servicio
	//Creamos la request y un Recorder que "guardara" la respuesta de la api
	r, _ := http.NewRequest("GET", "/sendjson", nil)
	w := httptest.NewRecorder()
	//Podemos ya probar el mock server con la request
	http.DefaultServeMux.ServeHTTP(w, r)

	var u struct {
		Name  string
		Email string
	}

	//Tomamos la respuesta, que será un json, y creamos el objeto u
	if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
		log.Println("ERROR:", err)
	}

	fmt.Println(u)
	// Output:
	// {Bill bill@ardanstudios.com}
}
