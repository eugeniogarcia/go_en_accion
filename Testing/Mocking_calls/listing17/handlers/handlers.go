// Package handlers provides the endpoints for the web service.
package handlers

import (
	"encoding/json"
	"net/http"
)

// Routes sets the routes for the web service.
// Define todas las rutas que vamos a mockear
func Routes() {
	http.HandleFunc("/sendjson", SendJSON)
}

// SendJSON returns a simple JSON document.
// Define un handler
func SendJSON(rw http.ResponseWriter, r *http.Request) {
	u := struct {
		Name  string
		Email string
	}{
		Name:  "Bill",
		Email: "bill@ardanstudios.com",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(&u)
}
