package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	// Registrar rutas
	mux.HandleFunc("/hola", HandleHola)
	mux.HandleFunc("/adios", HandleAdios)
	mux.HandleFunc("/status", HandleStatus)

	// El orden importa: los primeros middlewares se ejecutan primero
	var handler http.Handler = mux
	handler = Logging(handler)            // se ejecuta 1º - registra entrada
	handler = CORS(handler)               // se ejecuta 2º - maneja CORS
	handler = Security("GOPHER")(handler) // se ejecuta 3º - valida seguridad
	handler = Timing(handler)             // se ejecuta 4º - mide tiempo

	// cualquier petición que trate el servidor le medimos el tiempo, validamos la seguridad, manejamos CORS, la registramos y finalmente la enviamos al mux para que se aplique la lógica de negocio específica de la ruta

	s := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	slog.Info("servidor iniciado", "address", s.Addr)
	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// ============================================
// HANDLERS
// ============================================

func HandleHola(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola Mundo!\n"))
}

func HandleAdios(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Adios Mundo!\n"))
}

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: OK\n")
}

// ============================================
// MIDDLEWARES - Cada uno maneja UNA responsabilidad
// ============================================

// Logging: registra cada request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Registrar detalles del request, pero no altera ni el request ni la response
		slog.Info("REQUEST",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)
		// Llamar al siguiente handler en la cadena
		next.ServeHTTP(w, r)
	})
}

// CORS: maneja headers CORS
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Agregar headers CORS a la respuesta
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Secret-Password")

		// Si es un preflight request (OPTIONS), respondemos aquí
		// Si se trata de un OPTIONS termina aqui el pipeline, no se llama al siguiente handler
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return // termina aquí el pipeline
		}
		// Continuar con el siguiente handler
		next.ServeHTTP(w, r)
	})
}

// Security: valida el password en headers
func Security(password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("You didn't give the secret password\n"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Timing: mide el tiempo de ejecución de cada request
func Timing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// se llama al siguiente handler en la cadena a "la mitad" del middleware
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		slog.Info("REQUEST_DURATION",
			"path", r.URL.Path,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// ============================================
// UTILIDAD: Encadenar múltiples middlewares
// ============================================

// Chain: aplica una lista de middlewares en orden (útil para rutas específicas)
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	// Los middlewares se aplican en ORDEN INVERSO para que se ejecuten en el orden deseado
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
