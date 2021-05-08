// This sample program demonstrates how to decode a JSON string.
//Trabajamos con un string que representa un Json (el stringify del json)
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Contact represents our JSON string.
type Contact struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Contact struct {
		Home string `json:"home"`
		Cell string `json:"cell"`
	} `json:"contact"`
}

// JSON contains a sample string to unmarshal.
var JSON = `{
	"name": "Gopher",
	"title": "programmer",
	"contact": {
		"home": "415.333.3333",
		"cell": "415.555.5555"
	}
}`

func main() {
	// Unmarshal the JSON string into our variable.
	var c Contact
	//Pasa el string a un []byte, y luego hace el unmarshall al type
	err := json.Unmarshal([]byte(JSON), &c)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	fmt.Println(c)
}
