package main

import (
	"log/slog"
	"net/http"
	"time"
)

func main() {
	terribleSecurity := ProveedorSeguridad("GOPHER")

	// creamos un multiplexador de rutas...
	mux := http.NewServeMux()

	// ...añadimos una ruta al mux usando un handler
	mux.Handle("/hola", terribleSecurity(Cronometra(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hola Mundo!\n"))
		}))))

	// ...añadimos una ruta al mux usando una función (la función tiene que tener la firma que se corresponde con un handler)
	mux.HandleFunc("/adios", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Adios Mundo!\n"))
	})
	//	mux = terribleSecurity(RequestTimer(mux))

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func Cronometra(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //devuelve un handler...
		start := time.Now()
		h.ServeHTTP(w, r) //...que ejecuta el handler que le pasamos
		end := time.Now()
		slog.Info("request time", "path", r.URL.Path, "duration", end.Sub(start))
	})
}

var securityMsg = []byte("You didn't give the secret password\n")

// Pasamos una configuracion para el middleware
func ProveedorSeguridad(password string) func(http.Handler) http.Handler {
	//y creamos un middleware...
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(securityMsg)
				return //abortamos la cadena de handlers
			}
			h.ServeHTTP(w, r) //continuamos con el siguiente handler
		})
	}
}
