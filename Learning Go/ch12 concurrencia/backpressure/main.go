package main

import (
	"errors"
	"net/http"
	"time"
)

type MiThrotle struct {
	ch chan struct{} //canal que usamos para limitar concurrencia
}

// Factoria que implementa un Mithrotle, que limita a limite ejecuciones concurrentes
func New(limite int) *MiThrotle {
	return &MiThrotle{
		ch: make(chan struct{}, limite), //crea un buffered channel con capacidad limite
	}
}

// metodo que implementa el procesamiento. Toma como argumento una función
func (pg *MiThrotle) Procesa(f func()) error {
	select {
	case pg.ch <- struct{}{}: //con cada llamada a Procesa añadimos un token al canal. Si el canal esta lleno, el case no se ejecuta con lo que conseguimos el efecto throtle que perseguimos
		f()     // ejecuta la función
		<-pg.ch //saca del canal un token, de modo que podamos hacer otra llamada
		return nil
	default:
		return errors.New("no hay capacidad disponible") //optamos por devolver un error si el canal esta lleno. Si quitaramos el default, la llamada se bloquearia hasta que hubiera capacidad
	}
}

func simulaTrabajo() string {
	time.Sleep(2 * time.Second)
	return "done"
}

func main() {
	pg := New(10)
	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		err := pg.Procesa(func() {
			w.Write([]byte(simulaTrabajo()))
		})
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Demasiadas peticiones, prueba mas tarde"))
		}
	})
	http.ListenAndServe(":8080", nil)
}
