package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func main() {
	// Creamos un contexto padre con timeout de 3 segundos
	ctx, cancelFuncParent := context.WithTimeout(context.Background(), 3*time.Second)
	// Aseguramos que almenos se llame a cancel una vez para liberar recursos
	defer cancelFuncParent()
	// Creamos un contexto hijo que puede ser cancelado con un motivo
	ctx, cancelFunc := context.WithCancelCause(ctx)
	// Aseguramos que almenos se llame a cancel una vez para liberar recursos
	defer cancelFunc(nil)

	// Canal para recibir resultados de las goroutines
	ch := make(chan string)

	// WaitGroup para controlar las goroutines que hay en ejecución
	var wg sync.WaitGroup
	// ... vamos a lanzar dos goroutines que hacen peticiones HTTP
	wg.Add(2)

	// Primera goroutine: hace peticiones a una URL que devuelve varios status
	go func() {
		// Aseguramos que al terminar la goroutine llamamos a Done en el WaitGroup, para que la cuenta de goroutines en ejecución se reduzca
		defer wg.Done()
		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/status/200,200,200,500")
			if err != nil {
				// Cancelamos indicando cual es el motivo de la cancelación
				cancelFunc(fmt.Errorf("in status goroutine: %w", err))
				return
			}
			if resp.StatusCode == http.StatusInternalServerError {
				// Cancelamos indicando cual es el motivo de la cancelación
				cancelFunc(errors.New("bad status"))
				return
			}
			ch <- "success from status"
			time.Sleep(1 * time.Second)
		}
	}()

	// Segunda goroutine: hace peticiones a una URL que tarda en responder
	go func() {
		// Aseguramos que al terminar la goroutine llamamos a Done en el WaitGroup, para que la cuenta de goroutines en ejecución se reduzca

		defer wg.Done()
		for {
			resp, err := makeRequest(ctx, "http://httpbin.org/delay/1")
			if err != nil {
				fmt.Println("in delay goroutine:", err)
				// Cancelamos indicando cual es el motivo de la cancelación
				cancelFunc(fmt.Errorf("in delay goroutine: %w", err))
				return
			}
			ch <- "success from delay: " + resp.Header.Get("date")
		}
	}()

	// Bucle principal: escuchamos resultados de cada gorutina, y nos fijamos si se ha cancelado el contexto para terminar
loop:
	for {
		select {
		case s := <-ch:
			fmt.Println("in main:", s)
		case <-ctx.Done(): // entramos aquí si el contexto se ha cancelado
			fmt.Println("in main: cancelled with cause:", context.Cause(ctx), "err:", ctx.Err())
			break loop
		}
	}

	// Esperamos a que todas las goroutines terminen
	wg.Wait()
	fmt.Println("context cause:", context.Cause(ctx))
}
