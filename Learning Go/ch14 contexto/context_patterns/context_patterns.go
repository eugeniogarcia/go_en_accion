package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// main shows how to create a context and pass it to a function
func main() {
	ctx := context.Background()
	result, err := logic(ctx, "a string")
	fmt.Println(result, err)
}

// logic shows the parameters for functions that pass or use the context
func logic(ctx context.Context, info string) (string, error) {
	// do some interesting stuff here
	return "", nil
}

// Middleware. Toma un handler y devuelve un handler
func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Paso 1. Obtenemos de la request el contexto
		ctx := req.Context()
		// Paso 2. Hacemos algo con el contexto
		// Paso 3. Asociamos a la request este contexto
		req = req.WithContext(ctx)

		// Paso 4. Llamamos al siguiente handler de la cadena, pero a partir de este punto el contexto ya esta enriquecido
		handler.ServeHTTP(rw, req)
	})
}

// Handler. Ejemplo de un handler que procesa una request HTTP
func handler(rw http.ResponseWriter, req *http.Request) {
	// Paso 1. Obtenemos de la request el contexto
	ctx := req.Context()

	// Paso 2. Aplicamos la logica de negocio, usando el contexto cuando proceda
	err := req.ParseForm()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	data := req.FormValue("data")
	result, err := logic(ctx, data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write([]byte(result))
}

type ServiceCaller struct {
	client *http.Client
}

// Demostramos como propagar el contexto a otros servicios de modo que tenemos una cadena de servicios que comparten el mismo contexto
func (sc ServiceCaller) callAnotherService(ctx context.Context, data string) (string, error) {
	// Paso 1. Creamos la request HTTP, asociando el contexto (no creamos un contexto nuevo, sino que usamos el que nos pasan)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"http://example.com?data="+data, nil)
	if err != nil {
		return "", err
	}
	resp, err := sc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status code %d",
			resp.StatusCode)
	}
	// Procesa la respuesta
	id, err := processResponse(resp.Body)
	return id, err
}

// processResponse is a placeholder function for processing the body of an *http.Response
func processResponse(body io.ReadCloser) (string, error) {
	return "", nil
}
