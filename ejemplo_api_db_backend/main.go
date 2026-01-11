package main

import (
	"log"
	"os"
	"runners-postgresql/config"
	"runners-postgresql/server"

	_ "github.com/lib/pq"
)

func main() {
	log.Println("Starting Runners App")

	log.Println("Initializing configuration")
	// recuperamos la configuración
	config := config.InitConfig(getConfigFileName())

	log.Println("Initializing database")
	// inicializamos la base de datos
	dbHandler := server.InitDatabase(config)

	log.Println("Initializing Prometheus")
	// inicializamos Prometheus
	go server.InitPrometheus()

	log.Println("Initializig HTTP sever")
	// inicializamos el servidor HTTP donde exponemos los diferentes recursos y  las apis asociadas a ellos. Pasamos la configuración y la base de datos
	httpServer := server.InitHttpServer(config, dbHandler)

	// arrancamos el servidor HTTP
	httpServer.Start()
}

func getConfigFileName() string {
	// buscamos la variable de entorno ENV
	env := os.Getenv("ENV")

	// si la variable de entorno ENV está definida, usamos un archivo de configuración específico
	if env != "" {
		return "runners-" + env
	}

	// si la variable de entorno ENV no está definida, usamos el archivo de configuración por defecto
	return "runners"
}
