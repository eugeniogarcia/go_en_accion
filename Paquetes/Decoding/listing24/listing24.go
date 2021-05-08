// This sample program demonstrates how to decode a JSON response
// using the json package and NewDecoder function.
//Toma un json y lo decodifica en un type
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type (
	// gResult maps to the result document received from the search.
	gResult struct {
		Nombre       string  `json:"name"`
		UnescapedURL string  `json:"unescapedUrl"`
		id           int     `json:"id"`
		Latitud      float64 `json:"latitude"`
		Longitud     float64 `json:"longitude"`
		Altura       float64 `json:"altitude"`
		Velocidad    float64 `json:"velocity"`
		Visibilidad  string  `json:"visibility"`
		Lat_Solar    float64 `json:"solar_lat"`
		Lon_Solar    float64 `json:"solar_lon"`
		Unidades     string  `json:"units"`
		Dia          float64 `json:"daynum"`
	}
)

func main() {
	uri := "https://api.wheretheiss.at/v1/satellites/25544"

	// Issue the search against Google.
	resp, err := http.Get(uri)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()

	// Decode the JSON response into our struct type.
	var gr gResult
	err = json.NewDecoder(resp.Body).Decode(&gr)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	fmt.Println(gr)

	// Marshal the struct type into a pretty print
	// version of the JSON document.
	pretty, err := json.MarshalIndent(gr, "", "    ")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	fmt.Println(string(pretty))
}
