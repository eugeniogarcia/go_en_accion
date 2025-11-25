package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	// Creamos una petición HTTP con el contexto proporcionado
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Realizamos la petición HTTP usando el cliente por defecto
	return http.DefaultClient.Do(req)
}

func main() {
	// Creamos un contexto que se puede cancelar
	ctx, cancelFunc := context.WithCancel(context.Background())
	// Aseguramos que se llame a cancelFunc al final de main. Siempre hay que llamar a la función de cancelación para liberar recursos.
	defer cancelFunc()

	// Canal para comunicar resultados entre goroutines y main
	ch := make(chan string)

	// Usamos un WaitGroup para esperar a que las goroutines terminen. Controlamos que haya dos goroutines.
	var wg sync.WaitGroup
	wg.Add(2)

	// Lanzamos la primera goroutine que hace peticiones a un endpoint que devuelve diferentes códigos de estado.
	go func() {
		// Aseguramos que se llame a wg.Done() al finalizar la goroutine, de modo que en wg tengamos el conteo correcto.
		defer wg.Done()

		// Implementa la lógica de esta primera goroutine. En este caso hacemos peticiones a un endpoint que devuelve diferentes códigos de estado (http status codes)
		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/status/200,200,200,500")
			if err != nil {
				fmt.Println("error in status goroutine:", err)
				// Si hay un error, cancelamos el contexto y salimos de la goroutine.
				cancelFunc()
				return
			}
			if resp.StatusCode == http.StatusInternalServerError {
				fmt.Println("bad status, exiting")
				// Si recibimos un 500, cancelamos el contexto y salimos de la goroutine.
				cancelFunc()
				return
			}
			// Si todo va bien, enviamos un mensaje al canal, de modo que main pueda recibirlo. También chequeamos si el contexto ha sido cancelado.
			select {
			case ch <- "success from status":
			case <-ctx.Done():
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Lanzamos la segunda goroutine que hace peticiones a un endpoint que retrasa la respuesta.
	go func() {
		// Aseguramos que se llame a wg.Done() al finalizar la goroutine, de modo que en wg tengamos el conteo correcto.
		defer wg.Done()

		// Implementa la lógica de esta segunda goroutine. En este caso hacemos peticiones a un endpoint que retrasa la respuesta.
		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/delay/1")
			if err != nil {
				fmt.Println("error in delay goroutine:", err)
				cancelFunc()
				return
			}
			select {
			case ch <- "success from delay: " + resp.Header.Get("date"):
			case <-ctx.Done():
			}
		}
	}()

	// Loop principal que recibe mensajes del canal o detecta la cancelación del contexto.
loop:
	for {
		select {
		case s := <-ch:
			fmt.Println("in main:", s)
		case <-ctx.Done():
			fmt.Println("in main: cancelled!")
			break loop
		}
	}
	wg.Wait()
}
