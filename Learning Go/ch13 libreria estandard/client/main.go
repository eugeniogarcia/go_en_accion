package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Puntero a un cliente HTTP con un timeout de 30 segundos
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// creamos una nueva solicitud HTTP GET con un contexto. Usamos el método Get, la URL del recurso y nil para indicar que no hay cuerpo en la solicitud.
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		panic(err)
	}
	// Completamos la request con un encabezado personalizado "X-My-Client"
	req.Header.Add("X-My-Client", "Learning Go")

	// Enviamos la solicitud HTTP y obtenemos la respuesta
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// Aseguramos el cierre del cuerpo de la respuesta al finalizar
	defer res.Body.Close()

	// Verificamos que el código de estado sea 200 OK
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status: got %v", res.Status))
	}

	// Imprimimos la cabecera Content-Type
	fmt.Println(res.Header.Get("Content-Type"))

	// Decodificamos el cuerpo de la respuesta JSON en una estructura Go
	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", data)
}
